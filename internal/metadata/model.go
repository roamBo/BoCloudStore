package metadata

type FileMetadata struct {
	FileID     string
	FileName   string
	TotalSize  int64
	ChunkCount int
	ChunkSize  int64
	Status     string
	UserID     string
	CreateAt   int64
	UpdateAt   int64
}

type ChunkMetadata struct {
	FileID      string
	ChunkID     string
	ETag        string
	Size        int64
	StoragePath string
}
