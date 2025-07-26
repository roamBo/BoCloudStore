package db

import (
	"context"
	"database/sql"
	"github.com/roamBo/BoCloudStore/internal/metadata"
)

type PostgresStore interface {
	InsertFile(ctx context.Context, file *metadata.FileMetadata) error
	InsertChunk(ctx context.Context, chunk *metadata.ChunkMetadata) error
	GetFile(ctx context.Context, fileID string) (*metadata.FileMetadata, error)
	UpdateFileStatus(ctx context.Context, fileID, status string) error
}
