package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/roamBo/BoCloudStore/internal/metadata"
	"go.uber.org/zap"
)

type RedisCache struct {
	client        *redis.Client
	logger        *zap.Logger
	defaultExpiry time.Duration
}

func NewRedisCache(client *redis.Client, logger *zap.Logger) *RedisCache {
	return &RedisCache{
		client:        client,
		logger:        logger,
		defaultExpiry: 24 * time.Hour,
	}
}

func (c *RedisCache) GetFileMetadata(ctx context.Context, fileID string) (*metadata.FileMetadata, error) {
	key := "file: metadata:" + fileID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		c.logger.Error("failed to get file metadata",
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return nil, err
	}
	var fileMetadata metadata.FileMetadata
	if err := json.Unmarshal(data, &fileMetadata); err != nil {
		c.logger.Error("failed to unmarshal file metadata",
			zap.String("fileID", fileID),
			zap.Error(err),
		)
	}
	return &fileMetadata, nil
}

func (c *RedisCache) SetFileMetadata(ctx context.Context, fileMeta *metadata.FileMetadata) error {
	key := "file: metadata:" + fileMeta.FileID
	data, err := json.Marshal(fileMeta)
	if err != nil {
		c.logger.Error("failed to marshal file metadata",
			zap.String("fileID", fileMeta.FileID),
			zap.Error(err),
		)
		return err
	}

	if err := c.client.Set(ctx, key, data, c.defaultExpiry).Err(); err != nil {
		c.logger.Error("failed to set file metadata",
			zap.String("fileID", fileMeta.FileID),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (c *RedisCache) DeleteFileMetadata(ctx context.Context, fileID string) error {
	key := "file: metadata:" + fileID
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("failed to delete file metadata",
			zap.String("fileID", fileID),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (c *RedisCache) BatchGet(ctx context.Context, fileIDs []string) (map[string]*metadata.FileMetadata, error) {
	keys := make([]string, len(fileIDs))
	for i, fileID := range fileIDs {
		keys[i] = "file: metadata:" + fileID
	}

	results, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		c.logger.Error("failed to get batch metadata", zap.Error(err))
		return nil, err
	}

	metaMap := make(map[string]*metadata.FileMetadata)
	for i, result := range results {
		if result == nil {
			continue
		}
		var fileMetadata metadata.FileMetadata
		if err := json.Unmarshal([]byte(result.(string)), &fileMetadata); err != nil {
			c.logger.Error("failed to unmarshal file metadata",
				zap.String("file_id", fileIDs[i]),
				zap.Error(err),
			)
			continue
		}
		metaMap[fileIDs[i]] = &fileMetadata
	}
	return metaMap, nil
}
