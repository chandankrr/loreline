package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/chandankrr/loreline/internal/database"
	"github.com/chandankrr/loreline/internal/model/user"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect user: %w", err)
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
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to collect user: %w", err)
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
		return nil,
			fmt.Errorf("failed to get user by id: %w", err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to collect user: %w", err)
	}

	return &userItem, nil
}

func (r *UserRepository) UpdateEmailVerified(
	ctx context.Context,
	db database.DBTX,
	userID uuid.UUID,
	verified bool,
) error {
	stmt := `
		UPDATE users
		SET
			email_verified = @email_verified
		WHERE
			id = @id
	`

	_, err := db.Exec(ctx, stmt, pgx.NamedArgs{
		"id":             userID,
		"email_verified": verified,
	})
	if err != nil {
		return fmt.Errorf("failed to update user email verification: %w", err)
	}

	return nil
}
