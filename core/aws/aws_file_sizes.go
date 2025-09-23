package aws

type fileSize int64

const (
	MaxFileSize         int64 = 500 * 1024 * 1024
	MaxChunkSize        int64 = 5 * 1024 * 1024
	MaxSingleUploadSize int64 = 5 * 1024 * 1024
)
