IF DB_ID(N'AttendanceHackday') IS NULL
BEGIN
    CREATE DATABASE AttendanceHackday;
END;
GO

USE AttendanceHackday;
GO

-- Delta seed for data added after Seed_DataBase.sql.
IF COL_LENGTH(N'dbo.AvatarConfigs', N'Effect') IS NULL
BEGIN
    ALTER TABLE dbo.AvatarConfigs ADD Effect nvarchar(200) NULL;
END;
GO

IF COL_LENGTH(N'dbo.ShopItems', N'ImagePath') IS NULL
BEGIN
    ALTER TABLE dbo.ShopItems ADD ImagePath nvarchar(300) NULL;
END;
GO

IF COL_LENGTH(N'dbo.ShopItems', N'Slot') IS NULL
BEGIN
    ALTER TABLE dbo.ShopItems ADD Slot nvarchar(50) NULL;
END;
GO

IF COL_LENGTH(N'dbo.Classrooms', N'StudentIDs') IS NOT NULL
BEGIN
    ALTER TABLE dbo.Classrooms DROP COLUMN StudentIDs;
END;
GO

IF OBJECT_ID(N'dbo.ManualCoinAdjustments', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.ManualCoinAdjustments (
        UserID nvarchar(100) NOT NULL PRIMARY KEY,
        Amount int NOT NULL
    );
END;
GO

IF OBJECT_ID(N'dbo.AvatarBaseImages', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.AvatarBaseImages (
        ID nvarchar(100) NOT NULL PRIMARY KEY,
        Label nvarchar(200) NOT NULL,
        ImagePath nvarchar(300) NOT NULL
    );
END;
GO

MERGE dbo.Users AS target
USING (VALUES
    (N'21', N'Peter', N'student', N'peter@example.com', N'$2a$10$JKh.VBpQzjvFdc9.SAvCUOm3/D95PqA/HjvQ3ixuwp7ySPspb4htO', N'classroom1')
) AS source (UserID, Name, Role, Email, PasswordHash, ClassroomID)
ON target.UserID = source.UserID
WHEN MATCHED THEN
    UPDATE SET
        Name = source.Name,
        Role = source.Role,
        Email = source.Email,
        PasswordHash = source.PasswordHash,
        ClassroomID = source.ClassroomID
WHEN NOT MATCHED THEN
    INSERT (UserID, Name, Role, Email, PasswordHash, ClassroomID)
    VALUES (source.UserID, source.Name, source.Role, source.Email, source.PasswordHash, source.ClassroomID);
GO

MERGE dbo.ClassroomStudents AS target
USING (VALUES
    (N'classroom1', N'21')
) AS source (ClassroomID, StudentID)
ON target.ClassroomID = source.ClassroomID
    AND target.StudentID = source.StudentID
WHEN NOT MATCHED THEN
    INSERT (ClassroomID, StudentID)
    VALUES (source.ClassroomID, source.StudentID);
GO

MERGE dbo.ShopItems AS target
USING (VALUES
    (N'aura_sparkle', N'Sparkle Aura', 9, N'A bright aura for a student who keeps showing up.', N'/static/images/cosmetics/aura_sparkle.png', N'effect'),
    (N'background_beach', N'Beach Background', 15, N'A soft shoreline with sand, sea, and gentle sky colors.', NULL, N'background'),
    (N'background_forest', N'Forest Background', 15, N'A calm woodland palette with mossy greens and warm light.', NULL, N'background'),
    (N'background_meadow', N'Meadow Background', 15, N'Pastel grass, tiny blooms, and a quiet afternoon feel.', NULL, N'background'),
    (N'background_mountain', N'Mountain Background', 15, N'Cool ridge colors with misty lavender shadows.', NULL, N'background'),
    (N'background_sky', N'Sky Background', 15, N'Airy blues and cloud-soft highlights for a clear day.', NULL, N'background'),
    (N'background_sunset', N'Sunset Background', 15, N'Peach, rose, and gold tones for a mellow evening glow.', NULL, N'background'),
    (N'cape_gold', N'Golden Cape', 12, N'A shiny cape for extra style.', N'/static/images/cosmetics/cape_gold.png', N'clothing'),
    (N'crown_flower', N'Flower Crown', 6, N'A leafy crown with bright classroom blooms.', N'/static/images/cosmetics/crown_flower.png', N'hair_style'),
    (N'glasses_rocket', N'Rocket Glasses', 10, N'A bold accessory for your avatar.', N'/static/images/cosmetics/glasses_rocket.png', N'accessory'),
    (N'hat_star', N'Star Hat', 5, N'A bright hat for a standout student.', N'/static/images/cosmetics/hat_star.png', N'hair_style'),
    (N'hat_wizard', N'Wizard Hat', 7, N'A tall purple hat for magical attendance streaks.', N'/static/images/cosmetics/hat_wizard.png', N'hair_style'),
    (N'headphones_gem', N'Gem Headphones', 9, N'Bright headphones with a gem on top.', N'/static/images/cosmetics/headphones_gem.png', N'accessory'),
    (N'hoodie_blue', N'Blue Hoodie', 7, N'A cozy hoodie for everyday questing.', N'/static/images/cosmetics/hoodie_blue.png', N'clothing'),
    (N'scarf_red', N'Red Scarf', 6, N'A bold scarf for chilly morning check-ins.', N'/static/images/cosmetics/scarf_red.png', N'clothing'),
    (N'shades_pixel', N'Pixel Shades', 8, N'Blocky shades with old-school cool.', N'/static/images/cosmetics/shades_pixel.png', N'accessory'),
    (N'trail_comet', N'Comet Trail', 11, N'A comet streak that follows your avatar.', N'/static/images/cosmetics/trail_comet.png', N'effect'),
    (N'trail_rainbow', N'Rainbow Trail', 8, N'A colorful trail effect for your avatar.', N'/static/images/cosmetics/trail_rainbow.png', N'effect')
) AS source (ID, Name, Price, Description, ImagePath, Slot)
ON target.ID = source.ID
WHEN MATCHED THEN
    UPDATE SET
        Name = source.Name,
        Price = source.Price,
        Description = source.Description,
        ImagePath = source.ImagePath,
        Slot = source.Slot
WHEN NOT MATCHED THEN
    INSERT (ID, Name, Price, Description, ImagePath, Slot)
    VALUES (source.ID, source.Name, source.Price, source.Description, source.ImagePath, source.Slot);
GO

MERGE dbo.AvatarBaseImages AS target
USING (VALUES
    (N'brainrot', N'BrainRot', N'/static/images/avatars/brainrot.png'),
    (N'd_money', N'D-Money', N'/static/images/avatars/d_money.png'),
    (N'funk_rapper', N'Funk Rapper', N'/static/images/avatars/funk_rapper.png'),
    (N'gerald', N'Gerald', N'/static/images/avatars/gerald.png'),
    (N'gopher', N'Gopher', N'/static/images/avatars/gopher.png'),
    (N'mike', N'Mike', N'/static/images/avatars/mike.png'),
    (N'milkman', N'Milk Man', N'/static/images/avatars/milkman.png'),
    (N'peter', N'Peter', N'/static/images/avatars/peter.png'),
    (N'salaryman', N'Salary Man', N'/static/images/avatars/salaryman.png')
) AS source (ID, Label, ImagePath)
ON target.ID = source.ID
WHEN MATCHED THEN
    UPDATE SET
        Label = source.Label,
        ImagePath = source.ImagePath
WHEN NOT MATCHED THEN
    INSERT (ID, Label, ImagePath)
    VALUES (source.ID, source.Label, source.ImagePath);
GO

MERGE dbo.OwnedShopItems AS target
USING (VALUES
    (N'21', N'headphones_gem'),
    (N'student1', N'glasses_rocket'),
    (N'student1', N'background_beach'),
    (N'student1', N'aura_sparkle')
) AS source (UserID, ShopItemID)
ON target.UserID = source.UserID
    AND target.ShopItemID = source.ShopItemID
WHEN NOT MATCHED THEN
    INSERT (UserID, ShopItemID)
    VALUES (source.UserID, source.ShopItemID);
GO

MERGE dbo.AvatarConfigs AS target
USING (VALUES
    (N'21', N'peter', N'', N'', N'headphones_gem', N''),
    (N'student1', N'funk_rapper', N'', N'', N'', N'aura_sparkle')
) AS source (UserID, Base, HairStyle, Clothing, Accessory, Effect)
ON target.UserID = source.UserID
WHEN MATCHED THEN
    UPDATE SET
        Base = source.Base,
        HairStyle = source.HairStyle,
        Clothing = source.Clothing,
        Accessory = source.Accessory,
        Effect = source.Effect
WHEN NOT MATCHED THEN
    INSERT (UserID, Base, HairStyle, Clothing, Accessory, Effect)
    VALUES (source.UserID, source.Base, source.HairStyle, source.Clothing, source.Accessory, source.Effect);
GO

MERGE dbo.ManualCoinAdjustments AS target
USING (VALUES
    (N'student1', 25)
) AS source (UserID, Amount)
ON target.UserID = source.UserID
WHEN MATCHED THEN
    UPDATE SET Amount = source.Amount
WHEN NOT MATCHED THEN
    INSERT (UserID, Amount)
    VALUES (source.UserID, source.Amount);
GO

MERGE dbo.Transactions AS target
USING (VALUES
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-06-30T15:35:29-06:00'), N'Attendance reward for 2026-06-30'),
    (N'student1', -10, CONVERT(datetimeoffset(7), N'2026-06-30T15:35:32-06:00'), N'Purchased Rocket Glasses'),
    (N'21', 1, CONVERT(datetimeoffset(7), N'2026-06-30T16:23:12-06:00'), N'Attendance reward for 2026-06-30'),
    (N'21', -9, CONVERT(datetimeoffset(7), N'2026-06-30T16:23:21-06:00'), N'Purchased Gem Headphones'),
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-07-01T16:08:04-06:00'), N'Attendance reward for 2026-07-01'),
    (N'student1', -15, CONVERT(datetimeoffset(7), N'2026-07-01T16:08:12-06:00'), N'Purchased Beach Background'),
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-07-02T14:39:02-06:00'), N'Attendance reward for 2026-07-02'),
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-07-02T09:07:28-06:00'), N'Attendance reward for 2026-07-02'),
    (N'student1', -9, CONVERT(datetimeoffset(7), N'2026-07-02T16:37:17-06:00'), N'Purchased Sparkle Aura'),
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-07-08T17:33:20-06:00'), N'Attendance reward for 2026-07-08'),
    (N'student1', 1, CONVERT(datetimeoffset(7), N'2026-07-13T15:10:46-06:00'), N'Attendance reward for 2026-07-13')
) AS source (UserID, Amount, Timestamp, Description)
ON target.UserID = source.UserID
    AND target.Amount = source.Amount
    AND target.Timestamp = source.Timestamp
    AND target.Description = source.Description
WHEN NOT MATCHED THEN
    INSERT (UserID, Amount, Timestamp, Description)
    VALUES (source.UserID, source.Amount, source.Timestamp, source.Description);
GO

MERGE dbo.AttendanceRecords AS target
USING (VALUES
    (N'student1', N'classroom1', N'["2026-06-09","2026-06-10","2026-06-15","2026-06-22","2026-06-24","2026-06-30","2026-07-01","2026-07-02","2026-07-08","2026-07-13"]', NULL),
    (N'21', N'classroom1', N'["2026-06-30"]', N'[]')
) AS source (UserID, ClassroomID, PresentDates, AbsentDates)
ON target.UserID = source.UserID
    AND target.ClassroomID = source.ClassroomID
WHEN MATCHED THEN
    UPDATE SET
        PresentDates = source.PresentDates,
        AbsentDates = source.AbsentDates
WHEN NOT MATCHED THEN
    INSERT (UserID, ClassroomID, PresentDates, AbsentDates)
    VALUES (source.UserID, source.ClassroomID, source.PresentDates, source.AbsentDates);
GO
