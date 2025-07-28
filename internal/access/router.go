package access

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/roamBo/BoCloudStore/internal/access/handlers"
	"github.com/roamBo/BoCloudStore/internal/access/middleware"
	"github.com/roamBo/BoCloudStore/internal/business/chunk_upload"
	"github.com/roamBo/BoCloudStore/internal/metadata/service"
	"github.com/roamBo/BoCloudStore/pkg/config"
	"go.uber.org/zap"
)

func SetupRouter(
	cfg *config.Config,
	minioClient *minio.Client,
	metadataSvc service.Service,
	chunkUploadSvc chunk_upload.Service,
	logger *zap.Logger,
) *gin.Engine {
	router := gin.Default()

	healthHandler := handlers.NewHealthHandler(minioClient, logger)

	router.GET("/health", healthHandler.HealthCheck)

	authMiddleware := middleware.JWTAuth(logger, cfg)

	uploadGroup := router.Group("/upload")
	uploadGroup.Use(authMiddleware)
	{
		uploadHandler := handlers.NewUploadHandler(chunkUploadSvc, metadataSvc, logger)
		uploadGroup.POST("/init", uploadHandler.InitUpload)                      // 初始化上传
		uploadGroup.POST("/:file_id/chunk/:chunk_id", uploadHandler.UploadChunk) // 上传分块
		uploadGroup.POST("/:file_id/merge", uploadHandler.MergeChunks)           // 合并分块
	}
	return router
}
