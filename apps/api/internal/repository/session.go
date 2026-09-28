package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/lib/utils/token"
	"github.com/chandankrr/loreline/internal/model/session"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository struct {
	server *server.Server
}

func NewSessionRepository(server *server.Server) *SessionRepository {
	return &SessionRepository{server: server}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
	db database.DBTX,
	userID uuid.UUID,
	ttl time.Duration,
	ipAddress, userAgent *string,
) (*session.Session, error) {
	stmt := `
		INSERT INTO
			sessions (
				user_id,
				token,
				expires_at,
				ip_address,
				user_agent
			)
		VALUES
			(
				@user_id,
				@token,
				@expires_at,
				@ip_address,
				@user_agent
			)
		RETURNING
		*
	`

	rawToken, err := token.Generate()
	if err != nil {
		return nil, err
	}

	tokenHash := token.Hash(rawToken)

	rows, err := db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":    userID,
		"token":      tokenHash,
		"expires_at": time.Now().UTC().Add(ttl),
		"ip_address": ipAddress,
		"user_agent": userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	sessionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[session.Session])
	if err != nil {
		return nil, fmt.Errorf("failed to collect session: %w", err)
	}

	// Return original token instead of hash token
	sessionItem.Token = rawToken

	return &sessionItem, nil
}

func (r *SessionRepository) GetSession(
	ctx context.Context,
	sessionToken string,
) (*session.Session, error) {
	stmt := `
		SELECT
			*
		FROM sessions
		WHERE
			token = @token
			AND revoked = false
	`

	tokenHash := token.Hash(sessionToken)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"token": tokenHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	sessionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[session.Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, fmt.Errorf("failed to collect session: %w", err)
	}

	// Return original token instead of hash token
	sessionItem.Token = sessionToken

	return &sessionItem, nil
}

func (r *SessionRepository) RevokeSession(
	ctx context.Context,
	db database.DBTX,
	sessionToken string,
) error {
	stmt := `
		UPDATE sessions
		SET
			revoked = true
		WHERE
			token = @token
			AND revoked = false
	`

	tokenHash := token.Hash(sessionToken)

	_, err := db.Exec(ctx, stmt, pgx.NamedArgs{
		"token": tokenHash,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

func (r *SessionRepository) RevokeAllUserSessions(
	ctx context.Context,
	userID uuid.UUID,
) error {
	stmt := `
		UPDATE sessions
		SET
			revoked = true
		WHERE
			user_id = @user_id
			AND revoked = false
	`

	_, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	return nil
}
