package access

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	"github.com/roamBo/BoCloudStore/internal/access/handlers"
)

func SetupRouter(minioClient *minio.Client, logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	healthHandler := handlers.NewHealthHandler(minioClient, logger)

	router.GET("/health", healthHandler.HealthCheck)

	return router
}
