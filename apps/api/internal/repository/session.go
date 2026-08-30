package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/chandankrr/loreline/internal/model/session"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SessionRepository struct {
	server *server.Server
}

func NewSessionRepository(server *server.Server) *SessionRepository {
	return &SessionRepository{server: server}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
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

	token, err := generateSecureToken()
	if err != nil {
		return nil, err
	}

	tokenHash := hashToken(token)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":    userID,
		"token":      tokenHash,
		"expires_at": time.Now().UTC().Add(ttl),
		"ip_address": ipAddress,
		"user_agent": userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create session query for user_id=%s: %w", userID.String(), err)
	}

	sessionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[session.Session])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:sessions for user_id=%s: %w", userID.String(), err)
	}

	// Return original token instead of hash token
	sessionItem.Token = token

	return &sessionItem, nil
}

func (r *SessionRepository) CreateSessionTx(
	ctx context.Context,
	tx pgx.Tx,
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

	token, err := generateSecureToken()
	if err != nil {
		return nil, err
	}

	tokenHash := hashToken(token)

	rows, err := tx.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":    userID,
		"token":      tokenHash,
		"expires_at": time.Now().UTC().Add(ttl),
		"ip_address": ipAddress,
		"user_agent": userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create session query for user_id=%s: %w", userID.String(), err)
	}

	sessionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[session.Session])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:sessions for user_id=%s: %w", userID.String(), err)
	}

	// Return original token instead of hash token
	sessionItem.Token = token

	return &sessionItem, nil
}

func (r *SessionRepository) GetSession(ctx context.Context, token string) (*session.Session, error) {
	stmt := `
		SELECT
			*
		FROM sessions
		WHERE
			token = @token
			AND revoked = false
	`

	tokenHash := hashToken(token)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"token": tokenHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get session by token query for token=%s: %w", token, err)
	}

	sessionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[session.Session])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:session for token=%s: %w", token, err)
	}

	// Return original token instead of hash token
	sessionItem.Token = token

	return &sessionItem, nil
}

func (r *SessionRepository) RevokeSession(ctx context.Context, token string) error {
	stmt := `
		UPDATE sessions
		SET
			revoked = true
		WHERE
			token = @token
			AND revoked = false
	`

	tokenHash := hashToken(token)

	_, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"token": tokenHash,
	})
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *SessionRepository) RevokeSessionTx(ctx context.Context, tx pgx.Tx, token string) error {
	stmt := `
		UPDATE sessions
		SET
			revoked = true
		WHERE
			token = @token
			AND revoked = false
	`

	tokenHash := hashToken(token)

	_, err := tx.Exec(ctx, stmt, pgx.NamedArgs{
		"token": tokenHash,
	})
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *SessionRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
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
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

// generateSecureToken generates a cryptographically secure session/refresh token
func generateSecureToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// hashToken hashes a token before database storage/lookup
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
