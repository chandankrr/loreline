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
