package book

import (
	"time"

	"github.com/chandankrr/loreline/internal/model"
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type BookStatus string

const (
	BookStatusProcessing BookStatus = "processing"
	BookStatusReady      BookStatus = "ready"
	BookStatusFailed     BookStatus = "failed"
)

type BookIndexingStatus string

const (
	BookIndexingStatusPending  BookIndexingStatus = "pending"
	BookIndexingStatusIndexing BookIndexingStatus = "indexing"
	BookIndexingStatusReady    BookIndexingStatus = "ready"
	BookIndexingStatusFailed   BookIndexingStatus = "failed"
)

type Book struct {
	model.Base
	UserID            uuid.UUID          `json:"userId" db:"user_id"`
	FolderID          *uuid.UUID         `json:"folderId" db:"folder_id"`
	Title             string             `json:"title" db:"title"`
	Author            *string            `json:"author" db:"author"`
	ObjectKey         string             `json:"objectKey" db:"object_key"`
	OriginalFilename  string             `json:"originalFilename" db:"original_filename"`
	MimeType          string             `json:"mimeType" db:"mime_type"`
	FileSize          int64              `json:"fileSize" db:"file_size"`
	PageCount         int                `json:"pageCount" db:"page_count"`
	Status            BookStatus         `json:"status" db:"status"`
	IndexingStatus    BookIndexingStatus `json:"indexingStatus" db:"indexing_status"`
	TotalChunks       int                `json:"totalChunks" db:"total_chunks"`
	IndexedChunks     int                `json:"indexedChunks" db:"indexed_chunks"`
	IndexingAttempts  int                `json:"indexingAttempts" db:"indexing_attempts"`
	IndexingError     *string            `json:"indexingError" db:"indexing_error"`
	IndexingUpdatedAt time.Time          `json:"indexingUpdatedAt" db:"indexing_updated_at"`
	LastPage          int                `json:"lastPage" db:"last_page"`
	Progress          float32            `json:"progress" db:"progress"`
	ErrorMessage      *string            `json:"errorMessage" db:"error_message"`
	LastOpenedAt      time.Time          `json:"lastOpenedAt" db:"last_opened_at"`
}

type BookChunk struct {
	model.Base
	BookID       uuid.UUID        `json:"bookId" db:"book_id"`
	UserID       uuid.UUID        `json:"userId" db:"user_id"`
	PageStart    int              `json:"pageStart" db:"page_start"`
	PageEnd      int              `json:"pageEnd" db:"page_end"`
	Content      string           `json:"content" db:"content"`
	SearchVector *string          `json:"searchVector,omitempty" db:"search_vector"`
	Embedding    *pgvector.Vector `json:"embedding,omitempty" db:"embedding"`
}
