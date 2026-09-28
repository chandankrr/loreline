package handler

import (
	"errors"
	"net/http"

	"github.com/chandankrr/loreline/internal/dto"
	"github.com/chandankrr/loreline/internal/errs"
	"github.com/chandankrr/loreline/internal/middleware"
	"github.com/chandankrr/loreline/internal/model/folder"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/chandankrr/loreline/internal/service"
	"github.com/labstack/echo/v4"
)

type FolderHandler struct {
	Handler
	folderService *service.FolderService
}

func NewFolderHandler(s *server.Server, folderService *service.FolderService) *FolderHandler {
	return &FolderHandler{
		Handler:       NewHandler(s),
		folderService: folderService,
	}
}

func (h *FolderHandler) CreateFolder(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.CreateFolderPayload) (*folder.Folder, error) {
			userID := middleware.GetUserID(c)

			createdFolder, err := h.folderService.CreateFolder(
				c.Request().Context(),
				userID,
				payload.Name,
				payload.ParentID,
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrParentFolderNotFound):
					return nil, errs.NewNotFoundError(
						"Parent folder doesn't exist",
						false, nil,
					)

				case errors.Is(err, service.ErrFolderAlreadyExists):
					code := "FOLDER_ALREADY_EXISTS"
					return nil, errs.NewConflictError(
						"Folder with this name already exists",
						false, &code,
					)

				default:
					return nil, err
				}
			}

			return createdFolder, nil
		},
		http.StatusCreated,
		&dto.CreateFolderPayload{},
	)(c)
}

func (h *FolderHandler) ListFolders(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.ListFoldersQuery) (*dto.ListFoldersResponse, error) {
			userID := middleware.GetUserID(c)

			folders, err := h.folderService.GetFolders(
				c.Request().Context(),
				userID,
				payload.ParentID,
			)
			if err != nil {
				return nil, err
			}

			return &dto.ListFoldersResponse{Folders: folders}, nil
		},
		http.StatusOK,
		&dto.ListFoldersQuery{},
	)(c)
}

func (h *FolderHandler) GetFolderTree(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.EmptyPayload) (*dto.GetFolderTreeResponse, error) {
			userID := middleware.GetUserID(c)

			tree, err := h.folderService.GetFolderTree(
				c.Request().Context(),
				userID,
			)
			if err != nil {
				return nil, err
			}

			return &dto.GetFolderTreeResponse{Tree: tree}, nil
		},
		http.StatusOK,
		&dto.EmptyPayload{},
	)(c)
}

func (h *FolderHandler) GetFolder(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.GetFolderPayload) (*folder.FolderSummary, error) {
			userID := middleware.GetUserID(c)

			folder, err := h.folderService.GetFolderSummary(
				c.Request().Context(),
				userID,
				payload.ID,
			)
			if err != nil {
				if errors.Is(err, service.ErrFolderNotFound) {
					return nil,
						errs.NewNotFoundError("Folder doesn't exist", false, nil)
				}

				return nil, err
			}

			return folder, nil
		},
		http.StatusOK,
		&dto.GetFolderPayload{},
	)(c)
}

func (h *FolderHandler) UpdateFolder(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.UpdateFolderPayload) (*folder.Folder, error) {
			userID := middleware.GetUserID(c)

			renamedFolder, err := h.folderService.RenameFolder(
				c.Request().Context(),
				userID,
				payload.ID,
				payload.Name,
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrFolderNotFound):
					return nil, errs.NewNotFoundError("Folder doesn't exist", false, nil)

				case errors.Is(err, service.ErrFolderAlreadyExists):
					code := "FOLDER_ALREADY_EXISTS"
					return nil, errs.NewConflictError(
						"Folder with this name already exists",
						false, &code,
					)

				default:
					return nil, err
				}

			}

			return renamedFolder, nil
		},
		http.StatusOK,
		&dto.UpdateFolderPayload{},
	)(c)
}

func (h *FolderHandler) DeleteFolder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *dto.DeleteFolderPayload) error {
			userID := middleware.GetUserID(c)

			err := h.folderService.DeleteFolder(
				c.Request().Context(),
				userID,
				payload.ID)
			if err != nil {
				if errors.Is(err, service.ErrFolderNotFound) {
					return errs.NewNotFoundError("Folder doesn't exist", false, nil)
				}

				return err
			}

			return nil
		},
		http.StatusNoContent,
		&dto.DeleteFolderPayload{},
	)(c)
}

func (h *FolderHandler) GetFolderBreadcrumb(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.GetFolderPayload) (*dto.FolderBreadcrumbResponse, error) {
			userID := middleware.GetUserID(c)

			breadcrumb, err := h.folderService.GetFolderBreadcrumb(
				c.Request().Context(),
				userID,
				payload.ID,
			)
			if err != nil {
				if errors.Is(err, service.ErrFolderNotFound) {
					return nil, errs.NewNotFoundError("Folder doesn't exist", false, nil)
				}

				return nil, err
			}

			return &dto.FolderBreadcrumbResponse{Breadcrumb: breadcrumb}, nil
		},
		http.StatusOK,
		&dto.GetFolderPayload{},
	)(c)
}

func (h *FolderHandler) MoveFolder(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *dto.MoveFolderPayload) (*folder.Folder, error) {
			userID := middleware.GetUserID(c)

			movedFolder, err := h.folderService.MoveFolder(
				c.Request().Context(),
				userID,
				payload.ID,
				payload.ParentID,
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrFolderNotFound):
					return nil, errs.NewNotFoundError("Folder doesn't exist", false, nil)

				case errors.Is(err, service.ErrParentFolderNotFound):
					return nil, errs.NewNotFoundError(
						"Destination folder doesn't exist",
						false, nil,
					)

				case errors.Is(err, service.ErrFolderCannotContainItself):
					return nil, errs.NewBadRequestError(
						"Folder cannot contain itself",
						false, nil, nil, nil,
					)

				case errors.Is(err, service.ErrFolderCannotMoveIntoDescendant):
					return nil, errs.NewBadRequestError(
						"Folder cannot be moved into its descendant",
						false, nil, nil, nil,
					)

				case errors.Is(err, service.ErrFolderAlreadyExists):
					code := "FOLDER_ALREADY_EXISTS"
					return nil, errs.NewConflictError(
						"Folder with this name already exists at destination",
						false, &code,
					)

				default:
					return nil, err
				}
			}

			return movedFolder, nil
		},
		http.StatusOK,
		&dto.MoveFolderPayload{},
	)(c)
}
