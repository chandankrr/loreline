package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/model/account"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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
		return nil, fmt.Errorf("failed to execute create account query for user_id=%s: %w", userID.String(), err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:accounts for user_id=%s: %w", userID.String(), err)
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
		return nil, fmt.Errorf("failed to execute get credential account query for user_id=%s: %w", userID.String(), err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:accounts for user_id=%s: %w", userID.String(), err)
	}

	return &accountItem, nil
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
		return nil, fmt.Errorf("failed to execute create account query for user_id=%s account_id=%s provider_id=%s: %w", userID.String(), accountID, providerID, err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:accounts for user_id=%s account_id=%s provider_id=%s: %w", userID, accountID, providerID, err)
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
		return nil, fmt.Errorf("failed to execute get account by provider_id and account_id query for provider_id=%s account_id=%s: %w", providerID, AccountID, err)
	}

	accountItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:accounts for provider_id=%s account_id=%s: %w", providerID, AccountID, err)
	}

	return &accountItem, nil
}

func (r *AccountRepository) UpdateOAuthTokens(
	ctx context.Context,
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

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"access_token":            accessToken,
		"refresh_token":           refreshToken,
		"id_token":                idToken,
		"access_token_expires_at": accessTokenExpiresAt,
		"id":                      id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	updatedAccount, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account.Account])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:accounts: %w", err)
	}

	return &updatedAccount, nil
}
