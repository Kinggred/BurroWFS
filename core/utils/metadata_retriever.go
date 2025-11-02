package utils

import (
	"burrowfs/core/logging"
	"bytes"
	"github.com/gabriel-vasile/mimetype"
	"io"
)

func GetFileMetadata(data io.ReadCloser) (newData io.ReadCloser, mime *mimetype.MIME, fileSize int64) {
	var logger = logging.Get("webdav/fileMetadata")
	buf, err := io.ReadAll(io.LimitReader(data, 4096))
	if err != nil {
		logger.Error("Failed to read data: ", err)
		return newData, nil, 0
	}
	mime = mimetype.Detect(buf)

	remainingData, err := io.ReadAll(data)
	if err != nil {
		logger.Error("Failed to read remaining data: ", err)
	}
	fullContents := append(buf, remainingData...)
	fileSize = int64(len(fullContents))
	data = io.NopCloser(bytes.NewReader(fullContents))

	logger.Debug("Detected MIME type: ", mime.String())
	return data, mime, fileSize
}
