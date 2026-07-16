package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

const startingStudentCoins = 10

var (
	ErrInvalidStudent          = errors.New("invalid student")
	ErrInvalidClassroom        = errors.New("invalid classroom")
	ErrAttendanceAlreadyMarked = errors.New("attendance already marked")
	ErrShopItemNotFound        = errors.New("shop item not found")
	ErrShopItemAlreadyOwned    = errors.New("shop item already owned")
	ErrInsufficientCoins       = errors.New("insufficient coins")
)

// LoadStudentDashboardState loads the shared student data, classroom schedule,
// and recurring assignment templates needed by the dashboard.
func (s *SQLStore) LoadStudentDashboardState(ctx context.Context, user domain.User) (domain.StudentState, error) {
	state, err := s.loadStudentPageState(ctx, user, true, false)
	if err != nil {
		return domain.StudentState{}, err
	}
	if err := s.loadWeeklyAssignmentTemplates(ctx, &state); err != nil {
		return domain.StudentState{}, err
	}
	return state, nil
}

// LoadStudentAttendanceState loads schedule data needed to calculate an
// attendance reward without fetching dashboard-only assignment templates.
func (s *SQLStore) LoadStudentAttendanceState(ctx context.Context, user domain.User) (domain.StudentState, error) {
	return s.loadStudentPageState(ctx, user, true, false)
}

// LoadStudentShopState loads the shared student data plus the SQL shop catalog.
func (s *SQLStore) LoadStudentShopState(ctx context.Context, user domain.User) (domain.StudentState, error) {
	return s.loadStudentPageState(ctx, user, false, true)
}

// LoadStudentAvatarState loads only the shared data used by the avatar page.
func (s *SQLStore) LoadStudentAvatarState(ctx context.Context, user domain.User) (domain.StudentState, error) {
	return s.loadStudentPageState(ctx, user, false, false)
}

// loadStudentPageState keeps page reads focused while loading balance,
// attendance, avatar, and ownership data shared by the student shell.
func (s *SQLStore) loadStudentPageState(ctx context.Context, user domain.User, includeSchedules, includeShop bool) (domain.StudentState, error) {
	if user.Role != "student" || user.ClassroomID == "" {
		return domain.StudentState{}, ErrInvalidStudent
	}

	state := domain.StudentState{
		User: user,
	}
	if err := s.loadStudentCore(ctx, &state); err != nil {
		return domain.StudentState{}, err
	}
	if includeSchedules {
		if err := s.loadSchedules(ctx, &state); err != nil {
			return domain.StudentState{}, err
		}
	}
	if includeShop {
		if err := s.loadShopItems(ctx, &state); err != nil {
			return domain.StudentState{}, err
		}
	}
	if err := s.loadOwnedShopItems(ctx, &state); err != nil {
		return domain.StudentState{}, err
	}
	return state, nil
}

// loadStudentCore combines the common balance, attendance, and avatar reads
// into one database round trip for every student page.
func (s *SQLStore) loadStudentCore(ctx context.Context, state *domain.StudentState) error {
	var (
		presentJSON string
		absentJSON  string
		hasAvatar   int
		config      domain.AvatarConfig
	)
	err := s.db.QueryRowContext(ctx, `
		SELECT
			CASE WHEN balance.Total < 0 THEN 0 ELSE balance.Total END,
			COALESCE(attendance.PresentDates, N'[]'),
			COALESCE(attendance.AbsentDates, N'[]'),
			CASE WHEN avatar.UserID IS NULL THEN 0 ELSE 1 END,
			COALESCE(avatar.Base, N''),
			COALESCE(avatar.HairStyle, N''),
			COALESCE(avatar.Clothing, N''),
			COALESCE(avatar.Accessory, N''),
			COALESCE(avatar.Effect, N'')
		FROM dbo.Users AS student
		CROSS APPLY (
			SELECT @p3
				+ COALESCE((SELECT Amount FROM dbo.ManualCoinAdjustments WHERE UserID = @p1), 0)
				+ COALESCE((SELECT SUM(Amount) FROM dbo.Transactions WHERE UserID = @p1), 0) AS Total
		) AS balance
		LEFT JOIN dbo.AttendanceRecords AS attendance
			ON attendance.UserID = student.UserID AND attendance.ClassroomID = @p2
		LEFT JOIN dbo.AvatarConfigs AS avatar ON avatar.UserID = student.UserID
		WHERE student.UserID = @p1;
	`, state.User.UserID, state.User.ClassroomID, startingStudentCoins).Scan(
		&state.CoinBalance,
		&presentJSON,
		&absentJSON,
		&hasAvatar,
		&config.Base,
		&config.HairStyle,
		&config.Clothing,
		&config.Accessory,
		&config.Effect,
	)
	if err != nil {
		return err
	}

	state.Attendance.UserID = state.User.UserID
	state.Attendance.ClassroomID = state.User.ClassroomID
	if err := json.Unmarshal([]byte(presentJSON), &state.Attendance.Present); err != nil {
		return fmt.Errorf("decode present dates: %w", err)
	}
	if err := json.Unmarshal([]byte(absentJSON), &state.Attendance.Absent); err != nil {
		return fmt.Errorf("decode absent dates: %w", err)
	}
	if hasAvatar == 1 {
		state.AvatarConfig = &config
	}
	return nil
}

func (s *SQLStore) loadSchedules(ctx context.Context, state *domain.StudentState) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ClassroomID, DayOfWeek, StartTime, EndTime, DoubleDay
		FROM dbo.Schedule WHERE ClassroomID = @p1 ORDER BY ScheduleID;
	`, state.User.ClassroomID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var schedule domain.Schedule
		if err := rows.Scan(&schedule.ClassroomID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.DoubleDay); err != nil {
			return err
		}
		state.Schedules = append(state.Schedules, schedule)
	}
	return rows.Err()
}

// loadWeeklyAssignmentTemplates reads the recurring classroom assignments used
// to build the student's current Sunday-through-Saturday dashboard calendar.
func (s *SQLStore) loadWeeklyAssignmentTemplates(ctx context.Context, state *domain.StudentState) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ClassroomID, DueWeekday, Subject, Title,
			CONVERT(varchar(5), DueTime, 108), DisplayOrder
		FROM dbo.WeeklyAssignmentTemplates
		WHERE ClassroomID = @p1
		ORDER BY DueWeekday, DisplayOrder, DueTime, Title;
	`, state.User.ClassroomID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment domain.WeeklyAssignmentTemplate
		if err := rows.Scan(
			&assignment.ClassroomID,
			&assignment.DueWeekday,
			&assignment.Subject,
			&assignment.Title,
			&assignment.DueTime,
			&assignment.DisplayOrder,
		); err != nil {
			return err
		}
		state.WeeklyAssignments = append(state.WeeklyAssignments, assignment)
	}
	return rows.Err()
}

func (s *SQLStore) loadShopItems(ctx context.Context, state *domain.StudentState) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ID, Name, Price, Description, COALESCE(ImagePath, N''), COALESCE(Slot, N'')
		FROM dbo.ShopItems ORDER BY ID;
	`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var item domain.ShopItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.ImagePath, &item.Slot); err != nil {
			rows.Close()
			return err
		}
		state.ShopItems = append(state.ShopItems, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return nil
}

func (s *SQLStore) loadOwnedShopItems(ctx context.Context, state *domain.StudentState) error {
	ownedRows, err := s.db.QueryContext(ctx, `
		SELECT ShopItemID FROM dbo.OwnedShopItems WHERE UserID = @p1 ORDER BY ShopItemID;
	`, state.User.UserID)
	if err != nil {
		return err
	}
	defer ownedRows.Close()
	for ownedRows.Next() {
		var itemID string
		if err := ownedRows.Scan(&itemID); err != nil {
			return err
		}
		state.OwnedShopItemIDs = append(state.OwnedShopItemIDs, itemID)
	}
	return ownedRows.Err()
}

// MarkAttendanceAndAwardCoins locks a student's attendance row so one date and
// its matching reward transaction are committed together at most once.
func (s *SQLStore) MarkAttendanceAndAwardCoins(ctx context.Context, userID, classroomID, date string, reward int, occurredAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var storedClassroomID, role string
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(ClassroomID, N''), Role FROM dbo.Users WITH (UPDLOCK, HOLDLOCK) WHERE UserID = @p1;
	`, userID).Scan(&storedClassroomID, &role)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrInvalidStudent
	}
	if err != nil {
		return err
	}
	if role != "student" {
		return ErrInvalidStudent
	}
	if storedClassroomID == "" || storedClassroomID != classroomID {
		return ErrInvalidClassroom
	}

	var presentJSON string
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(PresentDates, N'[]') FROM dbo.AttendanceRecords WITH (UPDLOCK, HOLDLOCK)
		WHERE UserID = @p1 AND ClassroomID = @p2;
	`, userID, classroomID).Scan(&presentJSON)
	recordExists := err == nil
	present := []string{}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err == nil {
		if err := json.Unmarshal([]byte(presentJSON), &present); err != nil {
			return fmt.Errorf("decode present dates: %w", err)
		}
	}
	for _, presentDate := range present {
		if presentDate == date {
			return ErrAttendanceAlreadyMarked
		}
	}
	present = append(present, date)
	encodedPresent, err := json.Marshal(present)
	if err != nil {
		return err
	}

	if !recordExists {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO dbo.AttendanceRecords (UserID, ClassroomID, PresentDates, AbsentDates)
			VALUES (@p1, @p2, @p3, N'[]');
		`, userID, classroomID, string(encodedPresent))
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE dbo.AttendanceRecords SET PresentDates = @p3 WHERE UserID = @p1 AND ClassroomID = @p2;
		`, userID, classroomID, string(encodedPresent))
	}
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO dbo.Transactions (UserID, Amount, Timestamp, Description)
		VALUES (@p1, @p2, @p3, @p4);
	`, userID, reward, occurredAt, fmt.Sprintf("Attendance reward for %s", date))
	if err != nil {
		return err
	}
	return tx.Commit()
}

// PurchaseShopItem checks ownership and balance under one serializable SQL
// transaction so the debit and ownership record cannot diverge.
func (s *SQLStore) PurchaseShopItem(ctx context.Context, userID, itemID string, occurredAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var itemName string
	var price int
	err = tx.QueryRowContext(ctx, `SELECT Name, Price FROM dbo.ShopItems WITH (HOLDLOCK) WHERE ID = @p1;`, itemID).Scan(&itemName, &price)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrShopItemNotFound
	}
	if err != nil {
		return err
	}
	var owned int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(1) FROM dbo.OwnedShopItems WITH (UPDLOCK, HOLDLOCK) WHERE UserID = @p1 AND ShopItemID = @p2;
	`, userID, itemID).Scan(&owned)
	if err != nil {
		return err
	}
	if owned > 0 {
		return ErrShopItemAlreadyOwned
	}

	var balance int
	err = tx.QueryRowContext(ctx, `
		SELECT @p2
			+ COALESCE((SELECT Amount FROM dbo.ManualCoinAdjustments WITH (HOLDLOCK) WHERE UserID = @p1), 0)
			+ COALESCE((SELECT SUM(Amount) FROM dbo.Transactions WITH (UPDLOCK, HOLDLOCK) WHERE UserID = @p1), 0);
	`, userID, startingStudentCoins).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < price {
		return ErrInsufficientCoins
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.Transactions (UserID, Amount, Timestamp, Description)
		VALUES (@p1, @p2, @p3, @p4);
	`, userID, -price, occurredAt, fmt.Sprintf("Purchased %s", itemName)); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dbo.OwnedShopItems (UserID, ShopItemID) VALUES (@p1, @p2);
	`, userID, itemID); err != nil {
		return err
	}
	return tx.Commit()
}

// SaveAvatarConfig replaces the student's saved avatar while preserving one row per user.
func (s *SQLStore) SaveAvatarConfig(ctx context.Context, userID string, config domain.AvatarConfig) error {
	_, err := s.db.ExecContext(ctx, `
		MERGE dbo.AvatarConfigs AS target
		USING (SELECT @p1 AS UserID) AS source ON target.UserID = source.UserID
		WHEN MATCHED THEN UPDATE SET Base = @p2, HairStyle = @p3, Clothing = @p4, Accessory = @p5, Effect = @p6
		WHEN NOT MATCHED THEN INSERT (UserID, Base, HairStyle, Clothing, Accessory, Effect)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	`, userID, config.Base, config.HairStyle, config.Clothing, config.Accessory, config.Effect)
	return err
}
