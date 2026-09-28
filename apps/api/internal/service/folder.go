package service

import (
	"context"
	"errors"
	"strings"

	"github.com/chandankrr/loreline/internal/logger"
	"github.com/chandankrr/loreline/internal/model/folder"
	"github.com/chandankrr/loreline/internal/repository"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/google/uuid"
)

var (
	ErrFolderNotFound                 = repository.ErrFolderNotFound
	ErrFolderAlreadyExists            = errors.New("folder already exists")
	ErrFolderCannotContainItself      = errors.New("folder cannot contain itself")
	ErrFolderCannotMoveIntoDescendant = errors.New("folder cannot be moved into its descendant")
	ErrParentFolderNotFound           = errors.New("parent folder not found")
)

type FolderService struct {
	server     *server.Server
	folderRepo *repository.FolderRepository
}

func NewFolderService(
	server *server.Server,
	folderRepo *repository.FolderRepository,
) *FolderService {
	return &FolderService{
		server:     server,
		folderRepo: folderRepo,
	}
}

func (s *FolderService) CreateFolder(
	ctx context.Context,
	userID uuid.UUID,
	name string,
	parentID *uuid.UUID,
) (*folder.Folder, error) {
	logger := logger.GetLogger(ctx)

	if parentID != nil {
		if _, err := s.folderRepo.GetFolderByID(ctx, userID, *parentID); err != nil {
			if errors.Is(err, ErrFolderNotFound) {
				logger.Error().
					Err(err).
					Str("parent_id", parentID.String()).
					Msg("parent folder not found")

				return nil, ErrParentFolderNotFound
			}

			logger.Error().
				Err(err).
				Str("parent_id", parentID.String()).
				Msg("failed to get parent folder")

			return nil, err
		}
	}

	siblings, err := s.folderRepo.GetFolders(ctx, userID, parentID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to get sibling folders")

		return nil, err
	}

	if hasSiblingWithName(siblings, name, nil) {
		logger.Warn().
			Str("name", name).
			Msg("folder already exists")

		return nil, ErrFolderAlreadyExists
	}

	createdFolder, err := s.folderRepo.CreateFolder(
		ctx,
		userID,
		name,
		parentID,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("name", name).
			Msg("failed to create folder")

		return nil, err
	}

	logger.Info().
		Str("event", "create_folder").
		Str("folder_id", createdFolder.ID.String()).
		Str("name", createdFolder.Name).
		Str("parent_id", func() string {
			if createdFolder.ParentID != nil {
				return createdFolder.ParentID.String()
			}
			return ""
		}()).
		Msg("folder created successfully")

	return createdFolder, nil
}

func (s *FolderService) GetFolders(
	ctx context.Context,
	userID uuid.UUID,
	parentID *uuid.UUID,
) ([]folder.FolderWithCounts, error) {
	logger := logger.GetLogger(ctx)

	folders, err := s.folderRepo.GetFolders(ctx, userID, parentID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("parent_id", func() string {
				if parentID != nil {
					return parentID.String()
				}
				return ""
			}()).
			Msg("failed to get folders")

		return nil, err
	}

	return folders, nil
}

func (s *FolderService) GetFolderTree(
	ctx context.Context,
	userID uuid.UUID,
) ([]folder.TreeRecord, error) {
	logger := logger.GetLogger(ctx)

	rows, err := s.folderRepo.GetAllFolders(ctx, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to get folders for tree")

		return nil, err
	}

	return buildFolderTree(rows), nil
}

func (s *FolderService) GetFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
) (*folder.Folder, error) {
	logger := logger.GetLogger(ctx)

	folder, err := s.folderRepo.GetFolderByID(ctx, userID, folderID)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			logger.Error().
				Err(err).
				Str("folder_id", folderID.String()).
				Msg("folder not found")

			return nil, err
		}

		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get folder")

		return nil, err
	}

	return folder, nil
}

func (s *FolderService) GetFolderSummary(
	ctx context.Context,
	userID uuid.UUID,
	folderID uuid.UUID,
) (*folder.FolderSummary, error) {
	logger := logger.GetLogger(ctx)

	f, err := s.folderRepo.GetFolderByID(ctx, userID, folderID)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			logger.Error().
				Err(err).
				Str("folder_id", folderID.String()).
				Msg("folder not found")

			return nil, err
		}

		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get folder for summary")

		return nil, err
	}

	subtree, err := s.folderRepo.GetSubtreeSummary(ctx, userID, folderID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get subtree summary")

		return nil, err
	}

	return &folder.FolderSummary{
		Folder:                *f,
		DescendantBookCount:   subtree.DescendantBookCount,
		DescendantFolderCount: subtree.DescendantFolderCount,
	}, nil
}

func (s *FolderService) RenameFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
	name string,
) (*folder.Folder, error) {
	logger := logger.GetLogger(ctx)

	existingFolder, err := s.folderRepo.GetFolderByID(ctx, userID, folderID)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			logger.Error().
				Err(err).
				Str("folder_id", folderID.String()).
				Msg("folder not found")

			return nil, err
		}

		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get folder for rename")

		return nil, err
	}

	siblings, err := s.folderRepo.GetFolders(ctx, userID, existingFolder.ParentID)
	if err != nil {
		logger.Error().
			Err(err).Str("folder_id", folderID.String()).
			Msg("failed to get sibling folders for rename")

		return nil, err
	}

	if hasSiblingWithName(siblings, name, &existingFolder.ID) {
		logger.Warn().
			Str("folder_id", folderID.String()).
			Str("name", name).
			Msg("folder rename rejected because name already exists")

		return nil, ErrFolderAlreadyExists
	}

	renamedFolder, err := s.folderRepo.RenameFolder(ctx, userID, folderID, name)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Str("name", name).
			Msg("failed to rename folder")

		return nil, err
	}

	logger.Info().
		Str("event", "rename_folder").
		Str("folder_id", folderID.String()).
		Str("name", renamedFolder.Name).
		Msg("folder renamed successfully")

	return renamedFolder, nil
}

func (s *FolderService) DeleteFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
) error {
	logger := logger.GetLogger(ctx)

	err := s.folderRepo.DeleteFolder(ctx, userID, folderID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to delete folder")

		return err
	}

	// TODO: Delete all related books from S3

	logger.Info().
		Str("event", "folder_deleted").
		Str("folder_id", folderID.String()).
		Msg("folder deleted successfully")

	return nil
}

func (s *FolderService) GetFolderBreadcrumb(
	ctx context.Context,
	userID, folderID uuid.UUID,
) ([]folder.Folder, error) {
	logger := logger.GetLogger(ctx)

	rows, err := s.folderRepo.GetAllFolders(ctx, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get folders for breadcrumb")

		return nil, err
	}

	breadcrumb := buildFolderBreadcrumb(rows, folderID)
	if len(breadcrumb) == 0 {
		logger.Warn().
			Str("folder_id", folderID.String()).
			Msg("folder not found while building breadcrumb")

		return nil, ErrFolderNotFound
	}

	return breadcrumb, nil
}

func (s *FolderService) MoveFolder(
	ctx context.Context,
	userID, folderID uuid.UUID,
	parentID *uuid.UUID,
) (*folder.Folder, error) {
	logger := logger.GetLogger(ctx)

	if parentID != nil && *parentID == folderID {
		logger.Warn().
			Str("folder_id", folderID.String()).
			Str("parent_id", parentID.String()).
			Msg("folder move rejected because folder cannot contain itself")

		return nil, ErrFolderCannotContainItself
	}

	existingFolder, err := s.folderRepo.GetFolderByID(ctx, userID, folderID)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			logger.Error().
				Err(err).
				Str("folder_id", folderID.String()).
				Msg("folder not found")

			return nil, err
		}

		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get folder for move")

		return nil, err
	}

	// If moving into another folder, the destination
	// isn't a descendant of the folder being moved.
	if parentID != nil {
		destination, err := s.GetFolderBreadcrumb(ctx, userID, *parentID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("folder_id", folderID.String()).
				Str("parent_id", parentID.String()).
				Msg("failed to get destination breadcrumb for folder move")

			return nil, err
		}

		for _, f := range destination {
			if f.ID == folderID {
				logger.Warn().
					Str("folder_id", folderID.String()).
					Str("parent_id", parentID.String()).
					Msg("folder move rejected because destination is a descendant")

				return nil, ErrFolderCannotMoveIntoDescendant
			}
		}
	}

	siblings, err := s.folderRepo.GetFolders(ctx, userID, parentID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Msg("failed to get destination folders for move")

		return nil, err
	}

	if hasSiblingWithName(siblings, existingFolder.Name, &existingFolder.ID) {
		logger.Warn().
			Str("folder_id", folderID.String()).
			Str("name", existingFolder.Name).
			Msg("folder move rejected because name already exists at destination")

		return nil, ErrFolderAlreadyExists
	}

	movedFolder, err := s.folderRepo.MoveFolder(ctx, userID, folderID, parentID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("folder_id", folderID.String()).
			Str("parent_id", func() string {
				if parentID != nil {
					return parentID.String()
				}
				return ""
			}()).
			Msg("failed to move folder")

		return nil, err
	}

	logger.Info().
		Str("event", "move_folder").
		Str("folder_id", folderID.String()).
		Msg("folder moved successfully")

	return movedFolder, nil
}

// hasSiblingWithName checks whether any folder in folders has the given
// name, excluding the folder whose ID matches excludedID (if provided).
func hasSiblingWithName(
	folders []folder.FolderWithCounts,
	name string,
	excludedID *uuid.UUID,
) bool {
	for _, f := range folders {
		if excludedID != nil && f.ID == *excludedID {
			continue
		}

		if strings.EqualFold(f.Name, name) {
			return true
		}

	}

	return false
}

func buildFolderTree(rows []folder.Folder) []folder.TreeRecord {
	childrenByParent := make(map[uuid.UUID][]folder.Folder, len(rows))

	var roots []folder.Folder

	for _, row := range rows {
		if row.ParentID != nil {
			childrenByParent[*row.ParentID] = append(childrenByParent[*row.ParentID], row)
		} else {
			roots = append(roots, row)
		}
	}

	return buildTreeChildren(roots, childrenByParent)
}

func buildTreeChildren(
	rows []folder.Folder,
	childrenByParent map[uuid.UUID][]folder.Folder,
) []folder.TreeRecord {
	tree := make([]folder.TreeRecord, 0, len(rows))

	for _, row := range rows {
		tree = append(tree, folder.TreeRecord{
			Folder:   row,
			Children: buildTreeChildren(childrenByParent[row.ID], childrenByParent),
		})
	}

	return tree
}

func buildFolderBreadcrumb(rows []folder.Folder, folderID uuid.UUID) []folder.Folder {
	byID := make(map[uuid.UUID]folder.Folder, len(rows))
	for _, row := range rows {
		byID[row.ID] = row
	}

	var trail []folder.Folder
	visited := make(map[uuid.UUID]bool)

	current, ok := byID[folderID]
	for ok && !visited[current.ID] {
		visited[current.ID] = true
		trail = append(trail, current)

		if current.ParentID == nil {
			break
		}
		current, ok = byID[*current.ParentID]
	}

	for i, j := 0, len(trail)-1; i < j; i, j = i+1, j-1 {
		trail[i], trail[j] = trail[j], trail[i]
	}

	return trail
}
