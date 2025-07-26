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
