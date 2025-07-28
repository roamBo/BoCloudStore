package chunk_upload

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/roamBo/BoCloudStore/internal/metadata/service"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/roamBo/BoCloudStore/internal/metadata"
	"github.com/roamBo/BoCloudStore/internal/storage"
	"github.com/roamBo/BoCloudStore/pkg/pool"
	"go.uber.org/zap"
)

type Service interface {
	UploadChunk(ctx context.Context, fileID string, chunkID int, data io.Reader, userID string) (*metadata.ChunkMetadata, error)
	MergeChunks(ctx context.Context, fileID string, userID string) error
}

type chunkUploadService struct {
	metadataSvc service.Service  //metadata service
	minioClient *minio.Client    //minio client
	bufferPool  *sync.Pool       //memory pool(for optimize performance)
	workerPool  *pool.WorkerPool //goroutines pool(for union chunk)
	logger      *zap.Logger
	chunkSize   int64 //chunk size
}

func NewService(
	metadataSvc service.Service,
	minioClient *minio.Client,
	workerPool *pool.WorkerPool,
	logger *zap.Logger,
	chunkSize int64,
) Service {
	return &chunkUploadService{
		metadataSvc: metadataSvc,
		minioClient: minioClient,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, chunkSize)
				return &buf
			},
		},
		workerPool: workerPool,
		logger:     logger,
		chunkSize:  chunkSize,
	}
}

func (s *chunkUploadService) UploadChunk(
	ctx context.Context,
	fileID string,
	chunkID int,
	data io.Reader,
	userID string,
) (*metadata.ChunkMetadata, error) {
	// 1. verify if file metadata exists (ensure initialized upload)
	fileMeta, err := s.metadataSvc.GetFileMetadata(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if fileMeta.UserID != userID {
		return nil, errors.New("permission denied: not file owner")
	}
	// 2.calculate chunk hash (for verification)
	hash := md5.New()
	tee := io.TeeReader(data, hash)

	// 3.store chunks to MinIO (path:{UserID}/{fileID}/chunk_{chunkID}
	storagePath := fmt.Sprintf("%s%s/chunk_%d", userID, fileID, chunkID)
	_, err = s.minioClient.PutObject(ctx, "cloudstor", storagePath, tee, -1, minio.PutObjectOptions{})
	if err != nil {
		s.logger.Error("failed to upload chunk to minio",
			zap.Error(err),
			zap.String("fileID", fileID),
			zap.Int("chunkID", chunkID))
		return nil, err
	}

	// 4. record chunks metadata
	chunkMeta := &metadata.ChunkMetadata{
		FileID:      fileID,
		ChunkID:     chunkID,
		ETag:        hex.EncodeToString(hash.Sum(nil)),
		StoragePath: storagePath,
	}
	if err := s.metadataSvc.SaveChunkMetadata(ctx, chunkMeta); err != nil {
		return nil, err
	}
	return chunkMeta, nil
}

func (s *chunkUploadService) MergeChunks(ctx context.Context, fileID string, userID string) error {
	// 1. verify file status and permissions
	fileMeta, err := s.metadataSvc.GetFileMetadata(ctx, fileID)
	if err != nil {
		return err
	}
	if fileMeta.Status == "merged" {
		return errors.New("file already merged")
	}
	if fileMeta.UserID != userID {
		return errors.New("permission denied: not file owner")
	}
	// 2. retrieve all chunk metadata (sorted by sequence number)
	chunks, err := s.metadataSvc.GetAll
	if err != nil {
		return err
	}
	if len(chunks) != fileMeta.ChunkCount {
		return errors.New("chunk count mismatch")
	}
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].ChunkID < chunks[j].ChunkID
	})
	// 3. merge partitions (using goroutines pool to read partitions in parallel and write them to the target file in sequence)
	destPath := fmt.Sprintf("%s%s/chunk_%d", userID, fileID, fileMeta.ChunkCount)
	// 4. update file status as merged

	return s.metadataSvc.UpdateFileStatus(ctx, fileID, "merged")
}
