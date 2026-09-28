package dto

import (
	"github.com/chandankrr/loreline/internal/model/folder"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateFolderPayload struct {
	Name     string     `json:"name" validate:"required,min=1,max=30"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"`
}

func (p *CreateFolderPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type ListFoldersQuery struct {
	ParentID *uuid.UUID `query:"parentId" validate:"omitempty,uuid"`
}

func (q *ListFoldersQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(q)
}

type ListFoldersResponse struct {
	Folders []folder.FolderWithCounts `json:"folders"`
}

type GetFolderPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetFolderPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateFolderPayload struct {
	ID   uuid.UUID `param:"id" validate:"required,uuid"`
	Name string    `json:"name" validate:"required,min=1,max=30"`
}

func (p *UpdateFolderPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeleteFolderPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteFolderPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type MoveFolderPayload struct {
	ID       uuid.UUID  `param:"id" validate:"required,uuid"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"`
}

func (p *MoveFolderPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetFolderTreeResponse struct {
	Tree []folder.TreeRecord `json:"tree"`
}

type FolderBreadcrumbResponse struct {
	Breadcrumb []folder.Folder `json:"breadcrumb"`
}
