package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/roamBo/BoCloudStore/internal/metadata"
	"time"
)

type PostgresStore interface {
	InsertFile(ctx context.Context, file *metadata.FileMetadata) error
	InsertChunk(ctx context.Context, chunk *metadata.ChunkMetadata) error
	GetFile(ctx context.Context, fileID string) (*metadata.FileMetadata, error)
	UpdateFileStatus(ctx context.Context, fileID, status string) error
}

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) PostgresStore {
	return &postgresStore{db: db}
}

func (p *postgresStore) InsertFile(ctx context.Context, file *metadata.FileMetadata) error {
	query := `
		INSERT INTO file_metadata (
			file_id, filename, total_size, chunk_count, 
			chunk_size, status, user_id, create_at, update_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	currentTime := time.Now().Unix()
	if file.CreateAt == 0 {
		file.CreateAt = currentTime
	}
	file.UpdateAt = currentTime

	_, err := p.db.ExecContext(
		ctx, query,
		file.FileID, file.FileName, file.TotalSize, file.ChunkCount,
		file.ChunkSize, file.Status, file.UserID, file.CreateAt, file.UpdateAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}
	return nil
}

func (p *postgresStore) InsertChunk(ctx context.Context, chunk *metadata.ChunkMetadata) error {
	query := `
		INSERT INTO chunk_metadata (
			file_id, chunk_id, etag, size, storage_path
		) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := p.db.ExecContext(
		ctx, query,
		chunk.FileID, chunk.ChunkID, chunk.ETag, chunk.Size, chunk.StoragePath,
	)

	if err != nil {
		return fmt.Errorf("failed to insert chunk metadata: %w", err)
	}
	return nil
}

func (p *postgresStore) GetFile(ctx context.Context, fileID string) (*metadata.FileMetadata, error) {
	query := `
		SELECT 
			file_id, filename, total_size, chunk_count, 
			chunk_size, status, user_id, create_at, update_at
		FROM file_metadata
		WHERE file_id = $1
	`

	row := p.db.QueryRowContext(ctx, query, fileID)

	file := &metadata.FileMetadata{}
	err := row.Scan(
		&file.FileID, &file.FileName, &file.TotalSize, &file.ChunkCount,
		&file.ChunkSize, &file.Status, &file.UserID, &file.CreateAt, &file.UpdateAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found: %s", fileID)
		}
		return nil, fmt.Errorf("failed to retrieve file metadata: %w", err)
	}
	return file, nil
}

func (p *postgresStore) UpdateFileStatus(ctx context.Context, fileID, status string) error {
	query := `
		UPDATE file_metadata
		SET status = $1, update_at = $2
		WHERE file_id = $3
	`

	updateAt := time.Now().Unix()
	result, err := p.db.ExecContext(ctx, query, status, updateAt, fileID)
	if err != nil {
		return fmt.Errorf("failed to update file status: %w", err)
	}

	// Check if record was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no file found with ID: %s", fileID)
	}
	return nil
}
