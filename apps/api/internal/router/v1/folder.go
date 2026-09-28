package v1

import (
	"github.com/chandankrr/loreline/internal/handler"
	"github.com/chandankrr/loreline/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerFolderRoutes(
	r *echo.Group,
	h *handler.FolderHandler,
	m *middleware.AuthMiddleware,
) {
	folders := r.Group("/folders", m.RequiredAuth)

	folders.POST("", h.CreateFolder)
	folders.GET("", h.ListFolders)
	folders.GET("/tree", h.GetFolderTree)
	folders.GET("/:id", h.GetFolder)
	folders.PATCH("/:id", h.UpdateFolder)
	folders.DELETE("/:id", h.DeleteFolder)
	folders.GET("/:id/breadcrumb", h.GetFolderBreadcrumb)
	folders.POST("/:id/move", h.MoveFolder)
}
