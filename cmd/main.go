package cmd

import (
	"github.com/roamBo/BoCloudStore/internal/access"
	"github.com/roamBo/BoCloudStore/internal/storage"
	"github.com/roamBo/BoCloudStore/pkg/config"
	"github.com/roamBo/BoCloudStore/pkg/utils"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger := utils.NewLogger(cfg.Env)
	minioClient, err := storage.NewMinioClient(storage.MinioConfig(cfg.Minio))

	if err != nil {
		logger.Fatal("Unable to initialize minio client", zap.Error(err))
	}

	router := access.SetupRouter(minioClient, logger)

	logger.Info("Starting server", zap.String("port", cfg.ServerPort))
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatal("Unable to start server", zap.Error(err))
	}
}
