package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/chandankrr/loreline/internal/model/folder"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrFolderNotFound = errors.New("folder not found")

type FolderRepository struct {
	server *server.Server
}

func NewFolderRepository(server *server.Server) *FolderRepository {
	return &FolderRepository{server: server}
}

func (r *FolderRepository) CreateFolder(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	parentID *uuid.UUID,
) (*folder.Folder, error) {
	stmt := `
        INSERT INTO
            folders (
                name,
                user_id,
                parent_id
            )
        VALUES
            (
                @name,
                @user_id,
                @parent_id
            )
        RETURNING
        *
    `

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":      name,
		"user_id":   userID,
		"parent_id": parentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	folderItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[folder.Folder])
	if err != nil {
		return nil, fmt.Errorf("failed to collect folder: %w", err)

	}

	return &folderItem, nil
}

func (r *FolderRepository) GetFolderByID(
	ctx context.Context,
	userID, folderID uuid.UUID,
) (*folder.Folder, error) {
	stmt := `
        SELECT
            *
        FROM folders
        WHERE
            id = @folder_id
            AND user_id = @user_id
    `

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"folder_id": folderID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get folder by id: %w", err)
	}

	folderItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[folder.Folder])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFolderNotFound
		}

		return nil, fmt.Errorf("failed to collect folder: %w", err)
	}

	return &folderItem, nil
}

func (r *FolderRepository) GetFolders(
	ctx context.Context,
	userID uuid.UUID,
	parentID *uuid.UUID,
) ([]folder.FolderWithCounts, error) {
	stmt := `
		SELECT
			f.*,
			(
				SELECT
					count(*)::int
				FROM books b
				WHERE
					b.folder_id = f.id
					AND b.user_id = @user_id
			) AS book_count,
			(
				SELECT
					count(*)::int
				FROM folders child
				WHERE
					child.parent_id = f.id
					AND child.user_id = @user_id
			) AS child_folder_count
		FROM folders f
		WHERE
			f.user_id = @user_id
			AND f.parent_id IS NOT DISTINCT FROM @parent_id
		ORDER BY
			f.name ASC,
			f.created_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"parent_id": parentID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	folders, err := pgx.CollectRows(rows, pgx.RowToStructByName[folder.FolderWithCounts])
	if err != nil {
		return nil, fmt.Errorf("failed to collect folders: %w", err)
	}

	return folders, nil
}

func (r *FolderRepository) GetAllFolders(
	ctx context.Context,
	userID uuid.UUID,
) ([]folder.Folder, error) {
	stmt := `
		SELECT
			*
		FROM folders
		WHERE
			user_id = @user_id
		ORDER BY
			name ASC,
			created_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all folders: %w", err)
	}

	folders, err := pgx.CollectRows(rows, pgx.RowToStructByName[folder.Folder])
	if err != nil {
		return nil, fmt.Errorf("failed to collect folders: %w", err)
	}

	return folders, nil
}

func (r *FolderRepository) GetSubtreeSummary(
	ctx context.Context,
	userID, folderID uuid.UUID,
) (*folder.SubtreeSummary, error) {
	stmt := `
		WITH RECURSIVE subtree AS (
			SELECT
				id
			FROM folders
			WHERE
				id = @folder_id
				AND user_id = @user_id
			UNION ALL
			SELECT
				child.id
			FROM folders child
			JOIN subtree parent ON child.parent_id = parent.id
			WHERE
				child.user_id = @user_id
		)
		SELECT
			(
				SELECT
					GREATEST(count(*) - 1, 0)::int
				FROM subtree
			) AS descendant_folder_count,
			(
				SELECT
					count(*)::int
				FROM books
				WHERE
					books.user_id = @user_id
					AND books.folder_id IN (SELECT id FROM subtree)
			) AS descendant_book_count,
			COALESCE(
				(
					SELECT
						array_agg(books.id)
					FROM books
					WHERE
						books.user_id = @user_id
						AND books.folder_id IN (SELECT id FROM subtree)
				), ARRAY[]::uuid[]
			) AS book_ids

	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"folder_id": folderID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get subtree summary: %w", err)
	}

	summary, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[folder.SubtreeSummary],
	)
	if err != nil {
		return nil, fmt.Errorf("failed to collect subtree summary: %w", err)
	}

	return &summary, nil
}

func (r *FolderRepository) RenameFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
	name string,
) (*folder.Folder, error) {
	stmt := `
        UPDATE folders
        SET
            name = @name
        WHERE
            id = @folder_id
            AND user_id = @user_id
        RETURNING
        *
    `

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":      name,
		"folder_id": folderID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to rename folder: %w", err)
	}

	renamedFolder, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[folder.Folder])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFolderNotFound
		}

		return nil, fmt.Errorf("failed to collect folder: %w", err)
	}

	return &renamedFolder, nil
}

func (r *FolderRepository) MoveFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
	parentID *uuid.UUID,
) (*folder.Folder, error) {
	stmt := `
		UPDATE folders
		SET
			parent_id = @parent_id
		WHERE
			id = @folder_id
			AND user_id = @user_id
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"parent_id": parentID,
		"folder_id": folderID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to move folder: %w", err)
	}

	folderItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[folder.Folder])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFolderNotFound
		}

		return nil, fmt.Errorf("failed to collect folder: %w", err)
	}

	return &folderItem, nil
}

func (r *FolderRepository) DeleteFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
) error {
	stmt := `
        DELETE
        FROM
            folders
        WHERE
            id = @folder_id
            AND user_id = @user_id
    `

	tag, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"folder_id": folderID,
		"user_id":   userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrFolderNotFound
	}

	return nil
}
