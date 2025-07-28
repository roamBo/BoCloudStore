package service

import (
	"context"
	"errors"
	"github.com/roamBo/BoCloudStore/internal/metadata"
	"github.com/roamBo/BoCloudStore/internal/metadata/cache"
	"github.com/roamBo/BoCloudStore/internal/metadata/db"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	CreateFileMetadata(ctx context.Context, file *metadata.FileMetadata) error
	SaveChunkMetadata(ctx context.Context, chunk *metadata.ChunkMetadata) error
	GetFileMetadata(ctx context.Context, fileID string) (*metadata.FileMetadata, error)
	UpdateFileStatus(ctx context.Context, fileID, status string) error
}

type metadataService struct {
	db     db.PostgresStore
	cache  cache.RedisCache
	logger *zap.Logger
}

func NewService(db db.PostgresStore, cache cache.RedisCache, logger *zap.Logger) Service {
	return &metadataService{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (m *metadataService) CreateFileMetadata(ctx context.Context, file *metadata.FileMetadata) error {
	// Set timestamps if not provided
	if file.CreateAt == 0 {
		file.CreateAt = time.Now().Unix()
	}
	file.UpdateAt = time.Now().Unix()

	// Insert into database
	if err := m.db.InsertFile(ctx, file); err != nil {
		m.logger.Error("failed to insert file metadata into database",
			zap.Error(err),
			zap.String("fileID", file.FileID))
		return errors.New("database operation failed")
	}

	// Cache file metadata (TTL: 1 hour)
	if err := m.cache.SetFileMetadata(ctx, file); err != nil {
		m.logger.Warn("failed to cache file metadata",
			zap.Error(err),
			zap.String("fileID", file.FileID))
		// Non-critical error, proceed without returning error
	}

	m.logger.Info("file metadata created successfully",
		zap.String("fileID", file.FileID),
		zap.String("status", file.Status))
	return nil
}
func (m *metadataService) SaveChunkMetadata(ctx context.Context, chunk *metadata.ChunkMetadata) error {
	if err := m.db.InsertChunk(ctx, chunk); err != nil {
		m.logger.Error("Failed to save chunk metadata to database",
			zap.Error(err),
			zap.String("fileID", chunk.FileID),
			zap.String("chunkID", chunk.ChunkID))
		return errors.New("database operation failed")
	}
	m.logger.Info("Chunk metadata saved successfully",
		zap.String("fileID", chunk.FileID),
		zap.String("chunkID", chunk.ChunkID))
	return nil
}
func (m *metadataService) GetFileMetadata(ctx context.Context, fileID string) (*metadata.FileMetadata, error) {
	// Try to get from cache first
	if file, err := m.cache.GetFileMetadata(ctx, fileID); err == nil {
		m.logger.Info("Retrieved file metadata from cache",
			zap.String("fileID", fileID))
		return file, nil
	}
	// Fallback to database if cache miss
	file, err := m.db.GetFile(ctx, fileID)
	if err != nil {
		m.logger.Error("Failed to retrieve file metadata from database",
			zap.Error(err),
			zap.String("fileID", fileID))
		return nil, errors.New("file not found")
	}

	if err := m.cache.SetFileMetadata(ctx, file); err != nil {
		m.logger.Warn("Failed to cache file metadata after retrieval",
			zap.Error(err),
			zap.String("fileID", fileID))
	}
	m.logger.Info("Retrieved file metadata from database",
		zap.String("fileID", fileID))
	return file, nil
}
func (m *metadataService) UpdateFileStatus(ctx context.Context, fileID, status string) error {
	if err := m.db.UpdateFileStatus(ctx, fileID, status); err != nil {
		m.logger.Error("Failed to update file status in database",
			zap.Error(err),
			zap.String("fileID", fileID),
			zap.String("status", status))
		return errors.New("database update failed")
	}

	// Invalidate cache to ensure consistency
	if err := m.cache.DeleteFileMetadata(ctx, fileID); err != nil {
		m.logger.Warn("Failed to invalidate cache after status update",
			zap.Error(err),
			zap.String("fileID", fileID))
	}

	m.logger.Info("File status updated successfully",
		zap.String("fileID", fileID),
		zap.String("newStatus", status))
	return nil
}
