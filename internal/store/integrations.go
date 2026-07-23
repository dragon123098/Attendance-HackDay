package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/PeterGrunig/Attendance-HackDay/internal/domain"
	"github.com/PeterGrunig/Attendance-HackDay/internal/integrations"
)

// CreateIntegrationConnection encrypts credentials before any database call,
// ensuring provider secrets can never be persisted in plaintext.
func (s *SQLStore) CreateIntegrationConnection(ctx context.Context, connection integrations.Connection, status string) (int64, error) {
	configuration := connection.Configuration
	if len(configuration) == 0 {
		configuration = json.RawMessage(`{}`)
	}
	if !json.Valid(configuration) {
		return 0, fmt.Errorf("%w: connection configuration must be valid JSON", integrations.ErrInvalidConfiguration)
	}

	var encrypted integrations.EncryptedCredentials
	if len(connection.Credentials) > 0 {
		if s.credentialCipher == nil {
			log.Printf("integration connection %q not stored: credential encryption unavailable", connection.DisplayName)
			return 0, integrations.ErrEncryptionUnavailable
		}
		var err error
		encrypted, err = s.credentialCipher.Encrypt(connection.Credentials)
		if err != nil {
			log.Printf("integration connection %q credential encryption failed: %v", connection.DisplayName, err)
			return 0, err
		}
	}

	var connectionID int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO IntegrationConnections
			(ProviderKind, ConnectionRole, DisplayName, Status, Configuration,
			 CredentialCiphertext, CredentialNonce, EncryptionVersion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING IntegrationConnectionID;
	`, connection.ProviderKind, connection.Role, connection.DisplayName, status, configuration,
		nullBytes(encrypted.Ciphertext), nullBytes(encrypted.Nonce), nullEncryptionVersion(encrypted.Version)).Scan(&connectionID)
	if err != nil {
		log.Printf("store integration connection %q: %v", connection.DisplayName, err)
		return 0, err
	}
	return connectionID, nil
}

// LoadIntegrationConnection decrypts credentials only for adapter runtime use;
// a missing or incorrect key leaves the stored ciphertext untouched.
func (s *SQLStore) LoadIntegrationConnection(ctx context.Context, connectionID int64) (integrations.Connection, string, error) {
	var connection integrations.Connection
	var status string
	var configuration, ciphertext, nonce []byte
	var version sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT IntegrationConnectionID, ProviderKind, ConnectionRole, DisplayName,
			Status, Configuration, CredentialCiphertext, CredentialNonce, EncryptionVersion
		FROM IntegrationConnections WHERE IntegrationConnectionID = $1;
	`, connectionID).Scan(&connection.ID, &connection.ProviderKind, &connection.Role,
		&connection.DisplayName, &status, &configuration, &ciphertext, &nonce, &version)
	if err != nil {
		return integrations.Connection{}, "", err
	}
	connection.Configuration = append(json.RawMessage(nil), configuration...)
	if len(ciphertext) == 0 {
		return connection, status, nil
	}
	if s.credentialCipher == nil {
		log.Printf("integration connection %d credentials unavailable: encryption key not configured", connectionID)
		return integrations.Connection{}, "", integrations.ErrEncryptionUnavailable
	}
	plaintext, err := s.credentialCipher.Decrypt(integrations.EncryptedCredentials{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Version:    int(version.Int64),
	})
	if err != nil {
		log.Printf("integration connection %d credential decryption failed: %v", connectionID, err)
		return integrations.Connection{}, "", err
	}
	connection.Credentials = append(json.RawMessage(nil), plaintext...)
	return connection, status, nil
}

func (s *SQLStore) UpsertExternalEntityMapping(ctx context.Context, mapping domain.ExternalEntityMapping) (int64, error) {
	var mappingID int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO ExternalEntityMappings
			(IntegrationConnectionID, EntityKind, ExternalID, LocalID, SISID, Active, LastSeenAt, UpdatedAt)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (IntegrationConnectionID, EntityKind, ExternalID) DO UPDATE SET
			LocalID = EXCLUDED.LocalID, SISID = EXCLUDED.SISID,
			Active = EXCLUDED.Active, LastSeenAt = EXCLUDED.LastSeenAt,
			UpdatedAt = CURRENT_TIMESTAMP
		RETURNING ExternalEntityMappingID;
	`, mapping.ConnectionID, mapping.EntityKind, mapping.ExternalID, mapping.LocalID,
		mapping.SISID, mapping.Active, mapping.LastSeenAt).Scan(&mappingID)
	return mappingID, err
}

// CreateAttendanceBatch stores an approved roster snapshot atomically so an
// exporter can never observe a batch without all of its student entries.
func (s *SQLStore) CreateAttendanceBatch(ctx context.Context, batch domain.AttendanceBatch) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var batchID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO AttendanceBatches
			(ClassroomID, AttendanceDate, Version, WorkflowState, ApprovedBy,
			 ApprovedAt, DestinationConnectionID)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6, $7)
		RETURNING AttendanceBatchID;
	`, batch.ClassroomID, batch.AttendanceDate, batch.Version, batch.WorkflowState,
		batch.ApprovedBy, batch.ApprovedAt, batch.DestinationConnectionID).Scan(&batchID)
	if err != nil {
		return 0, err
	}
	for _, entry := range batch.Entries {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO AttendanceBatchEntries
				(AttendanceBatchID, UserID, Status, AttendanceMarkID, DeliveryState,
				 ExternalRecordID, LastError)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), NULLIF($7, ''));
		`, batchID, entry.UserID, entry.Status, entry.AttendanceMarkID,
			entry.DeliveryState, entry.ExternalRecordID, entry.LastError); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return batchID, nil
}

func (s *SQLStore) RecordAttendanceDeliveryAttempt(ctx context.Context, attempt domain.AttendanceDeliveryAttempt) (int64, error) {
	var attemptID int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO AttendanceDeliveryAttempts
			(AttendanceBatchID, IntegrationConnectionID, IdempotencyKey, AttemptNumber,
			 State, ResponseCode, ErrorMessage, StartedAt, CompletedAt)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, $9)
		RETURNING AttendanceDeliveryAttemptID;
	`, attempt.BatchID, attempt.ConnectionID, attempt.IdempotencyKey,
		attempt.AttemptNumber, attempt.State, attempt.ResponseCode,
		attempt.ErrorMessage, attempt.StartedAt, attempt.CompletedAt).Scan(&attemptID)
	return attemptID, err
}

func (s *SQLStore) AppendIntegrationAuditEvent(ctx context.Context, event domain.IntegrationAuditEvent) (int64, error) {
	metadata := event.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	if !json.Valid(metadata) {
		return 0, fmt.Errorf("%w: audit metadata must be valid JSON", integrations.ErrInvalidConfiguration)
	}
	var eventID int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO IntegrationAuditEvents
			(ActorUserID, EventType, EntityType, EntityID, Metadata, OccurredAt)
		VALUES (NULLIF($1, ''), $2, $3, $4, $5, $6)
		RETURNING IntegrationAuditEventID;
	`, event.ActorUserID, event.EventType, event.EntityType, event.EntityID,
		metadata, event.OccurredAt).Scan(&eventID)
	return eventID, err
}

func nullBytes(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return value
}

func nullEncryptionVersion(version int) any {
	if version == 0 {
		return nil
	}
	return version
}
