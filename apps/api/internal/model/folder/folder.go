package folder

import (
	"github.com/chandankrr/loreline/internal/model"
	"github.com/google/uuid"
)

type Folder struct {
	model.Base
	Name     string     `json:"name" db:"name"`
	UserID   uuid.UUID  `json:"userId" db:"user_id"`
	ParentID *uuid.UUID `json:"parentId" db:"parent_id"`
}

type FolderWithCounts struct {
	Folder
	BookCount        int `json:"bookCount" db:"book_count"`
	ChildFolderCount int `json:"childFolderCount" db:"child_folder_count"`
}

type SubtreeSummary struct {
	BookIDs               []uuid.UUID `json:"bookIds" db:"book_ids"`
	DescendantBookCount   int         `json:"descendantBookCount" db:"descendant_book_count"`
	DescendantFolderCount int         `json:"descendantFolderCount" db:"descendant_folder_count"`
}

type TreeRecord struct {
	Folder
	Children []TreeRecord `json:"children"`
}

type FolderSummary struct {
	Folder
	DescendantBookCount   int `json:"descendantBookCount"`
	DescendantFolderCount int `json:"descedantFolderCount"`
}
