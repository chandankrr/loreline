package repository

import (
	"context"
	"fmt"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/model/user"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(server *server.Server) *UserRepository {
	return &UserRepository{server: server}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	db database.DBTX,
	name, email string,
	image *string,
	emailVerified bool,
) (*user.User, error) {
	stmt := `
		INSERT INTO
			users (
				name,
				email,
				image,
				email_verified
			)
		VALUES
			(
				@name,
				@email,
				@image,
				@email_verified
			)
		RETURNING
		*
	`

	rows, err := db.Query(ctx, stmt, pgx.NamedArgs{
		"name":           name,
		"email":          email,
		"image":          image,
		"email_verified": emailVerified,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query for name=%s email=%s: %w", name, email, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for name=%s email=%s: %w", name, email, err)
	}

	return &userItem, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM users
		WHERE
			email = @email
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by email query for email=%s: %w", email, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", email, err)
	}

	return &userItem, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM users
		WHERE
			id = @id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by id query for id=%s: %w", id.String(), err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for id=%s: %w", id.String(), err)
	}

	return &userItem, nil
}
