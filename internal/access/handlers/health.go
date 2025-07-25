package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HealthHandler struct {
	minioClient *minio.Client
	logger      *zap.Logger
}

func NewHealthHandler(minioClient *minio.Client, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		minioClient: minioClient,
		logger:      logger,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := h.minioClient.ListBuckets(ctx)
	if err != nil {
		h.logger.Error("Failed to list buckets", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "unhealthy",
			"error":  "storage service unavailable",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
