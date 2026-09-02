package repository

import (
	"context"
	"fmt"

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
