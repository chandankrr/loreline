package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/model/verification"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VerificationRepository struct {
	server *server.Server
}

func NewVerificationRepository(server *server.Server) *VerificationRepository {
	return &VerificationRepository{server: server}
}

func (r *VerificationRepository) Create(
	ctx context.Context,
	identifier, value string,
	verificationType verification.VerificationType,
	expiresAt time.Time,
) (*verification.Verification, error) {
	stmt := `
		INSERT INTO
			verification (
				identifier,
				type,
				value,
				expires_at
			)
		VALUES
			(
				@identifier,
				@type,
				@value,
				@expires_at
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"identifier": identifier,
		"type":       verificationType,
		"value":      value,
		"expires_at": expiresAt,
	})
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to execute create verification query for identifier=%s type=%s: %w",
				identifier,
				verificationType,
				err,
			)
	}

	verificationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[verification.Verification])
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to collect row from table:verification for identifier=%s type=%s: %w",
				identifier,
				verificationType,
				err,
			)
	}

	return &verificationItem, nil
}

func (r *VerificationRepository) GetLatestByIdentifierAndType(
	ctx context.Context,
	identifier string,
	verificationType verification.VerificationType,
) (*verification.Verification, error) {
	stmt := `
		SELECT
			*
		FROM verification
		WHERE
			identifier = @identifier
			AND type = @type
		ORDER BY created_at DESC
		LIMIT 1
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"identifier": identifier,
		"type":       verificationType,
	})
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to execute get latest verification query for identifier=%s type=%s: %w",
				identifier,
				verificationType,
				err,
			)
	}

	verificationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[verification.Verification])
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to collect row from table:verification for identifier=%s type=%s: %w",
				identifier,
				verificationType,
				err,
			)
	}

	return &verificationItem, nil
}

func (r *VerificationRepository) GetByValueAndType(
	ctx context.Context,
	value string,
	verificationType verification.VerificationType,
) (*verification.Verification, error) {
	stmt := `
		SELECT
			*
		FROM verification
		WHERE
			value = @value
			AND type = @type
		ORDER BY created_at DESC
		LIMIT 1
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"value": value,
		"type":  verificationType,
	})
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to execute get verification query for value=%s type=%s: %w",
				value,
				verificationType,
				err,
			)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[verification.Verification])
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to collect row from table:verification for value=%s type=%s: %w",
				value,
				verificationType,
				err,
			)
	}

	return &item, nil
}

func (r *VerificationRepository) DeleteByIdentifierAndType(
	ctx context.Context,
	db database.DBTX,
	identifier string,
	verificationType verification.VerificationType,
) error {
	stmt := `
		DELETE FROM verification
		WHERE
			identifier = @identifier
			AND type = @type
	`

	_, err := db.Exec(ctx, stmt, pgx.NamedArgs{
		"identifier": identifier,
		"type":       verificationType,
	})
	if err != nil {
		return fmt.Errorf(
			"failed to execute delete verification query for identifier=%s type=%s: %w",
			identifier,
			verificationType,
			err,
		)
	}

	return nil
}

func (r *VerificationRepository) DeleteExpired(
	ctx context.Context,
	db database.DBTX,
) error {
	stmt := `
		DELETE FROM verification
		WHERE
			expires_at < CURRENT_TIMESTAMP
	`

	_, err := db.Exec(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to execute delete expired verification query: %w", err)
	}

	return nil
}

func (r *VerificationRepository) IncrementAttempts(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	stmt := `
        UPDATE verification
        SET
			attempts = attempts + 1
        WHERE
			id = @id
        RETURNING
			attempts
    `

	var attempts int

	err := r.server.DB.Pool.QueryRow(
		ctx,
		stmt,
		pgx.NamedArgs{
			"id": id,
		},
	).Scan(&attempts)

	if err != nil {
		return 0, fmt.Errorf(
			"failed to execute increment verification attempts query: %w",
			err,
		)
	}

	return attempts, nil
}
