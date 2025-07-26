package cache

import (
	"context"
	"github.com/roamBo/BoCloudStore/internal/metadata"
)

type MetadataCache interface {
	//get single file meta data
	GetFileMetadata(ctx context.Context, fileID string) (*metadata.FileMetadata, error)
	SetFileMetadata(ctx context.Context, fileMeta *metadata.FileMetadata) error
	DeleteFileMetadata(ctx context.Context, fileID string) error
	BatchGet(ctx context.Context, fileIDs []string) (map[string]*metadata.FileMetadata, error)
}
