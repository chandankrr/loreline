package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/model/account"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountRepository struct {
	server *server.Server
}

func NewAccountRepository(server *server.Server) *AccountRepository {
	return &AccountRepository{server: server}
}

func (r *AccountRepository) CreateCredentialAccount(
	ctx context.Context,
	db database.DBTX,
	userID uuid.UUID,
	passwordHash string,
) (*account.Account, error) {
	stmt := `
		INSERT INTO
			accounts (
				account_id,
				provider_id,
				user_id,
				password
			)
		VALUES
			(
				@account_id,
				@provider_id,
				@user_id,
				@password
			)
		RETURNING
		*
	`

	rows, err := db.Query(ctx, stmt, pgx.NamedArgs{
		"account_id":  userID.String(),
		"provider_id": account.ProviderCredential,
		"user_id":     userID,
		"password":    passwordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credential account: %w", err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect account: %w", err)
	}

	return &accountItem, nil
}

func (r *AccountRepository) GetCredentialAccount(
	ctx context.Context,
	userID uuid.UUID,
) (*account.Account, error) {
	stmt := `
		SELECT
			*
		FROM accounts
		WHERE
			user_id = @user_id
			AND provider_id = 'credential'
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get credential account: %w", err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, fmt.Errorf("failed to collect account: %w", err)
	}

	return &accountItem, nil
}

func (r *AccountRepository) UpdatePassword(
	ctx context.Context,
	db database.DBTX,
	userID uuid.UUID,
	passwordHash string,
) error {
	stmt := `
		UPDATE accounts
		SET
			password = @password
		WHERE
			user_id = @user_id
			AND provider_id = 'credential'
	`

	_, err := db.Exec(ctx, stmt, pgx.NamedArgs{
		"user_id":  userID,
		"password": passwordHash,
	})
	if err != nil {
		return fmt.Errorf("failed to update account password: %w", err)
	}

	return nil
}

func (r *AccountRepository) CreateOAuthAccount(
	ctx context.Context,
	db database.DBTX,
	userID uuid.UUID,
	providerID, accountID string,
	accessToken, refreshToken, idToken *string,
	accessTokenExpiresAt *time.Time,
	scope *string,
) (*account.Account, error) {
	stmt := `
		INSERT INTO
			accounts (
				account_id,
				provider_id,
				user_id,
				access_token,
				refresh_token,
				id_token,
				access_token_expires_at,
				scope
			)
		VALUES
			(
				@account_id,
				@provider_id,
				@user_id,
				@access_token,
				@refresh_token,
				@id_token,
				@access_token_expires_at,
				@scope
			)
		RETURNING
		*
	`

	rows, err := db.Query(ctx, stmt, pgx.NamedArgs{
		"account_id":              accountID,
		"provider_id":             providerID,
		"user_id":                 userID,
		"access_token":            accessToken,
		"refresh_token":           refreshToken,
		"id_token":                idToken,
		"access_token_expires_at": accessTokenExpiresAt,
		"scope":                   scope,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth account: %w", err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect account: %w", err)
	}

	return &accountItem, nil
}

func (r *AccountRepository) GetByProviderAndAccountID(
	ctx context.Context,
	providerID, AccountID string,
) (*account.Account, error) {
	stmt := `
		SELECT
			*
		FROM accounts
		WHERE
			provider_id = @provider_id
			AND account_id = @account_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"provider_id": providerID,
		"account_id":  AccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get account by provider and account id: %w", err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, fmt.Errorf("failed to collect account: %w", err)
	}

	return &accountItem, nil
}

func (r *AccountRepository) UpdateOAuthTokens(
	ctx context.Context,
	db database.DBTX,
	id uuid.UUID,
	accessToken, refreshToken, idToken *string,
	accessTokenExpiresAt *time.Time,
) (*account.Account, error) {
	stmt := `
		UPDATE accounts
		SET
			access_token = @access_token,
			refresh_token = @refresh_token,
			id_token = @id_token,
			access_token_expires_at = @access_token_expires_at
		WHERE
			id = @id
		RETURNING
		*
	`

	rows, err := db.Query(ctx, stmt, pgx.NamedArgs{
		"access_token":            accessToken,
		"refresh_token":           refreshToken,
		"id_token":                idToken,
		"access_token_expires_at": accessTokenExpiresAt,
		"id":                      id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update oauth tokens: %w", err)
	}

	updatedAccount, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}

		return nil, fmt.Errorf("failed to collect account: %w", err)
	}

	return &updatedAccount, nil
}
