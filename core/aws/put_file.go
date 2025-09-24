package aws

import (
	"burrowfs/core/config"
	"burrowfs/core/db/models"
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func multiPartUpload(ctx context.Context, s3Client *s3.Client, file *models.File, data io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createResp, err := s3Client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      &config.CONFIG.AWSBucket,
		Key:         &file.S3Key,
		ContentType: &file.ContentType,
	})

	if err != nil {
		return "", err
	}

	uploadID := createResp.UploadId
	var completedParts []types.CompletedPart
	partNumber := int32(0)
	buffer := make([]byte, MaxChunkSize)

	for {
		n, err := io.ReadFull(data, buffer)
		if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
			abortInput := &s3.AbortMultipartUploadInput{
				Bucket:   &config.CONFIG.AWSBucket,
				Key:      &file.S3Key,
				UploadId: uploadID,
			}
			_, err := s3Client.AbortMultipartUpload(ctx, abortInput)
			if err != nil {
				return "", err
			}
			return "", err
		}
		if n == 0 {
			break
		}

		part := buffer[:n]
		uploadPartResp, err := s3Client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:     &config.CONFIG.AWSBucket,
			Key:        &file.S3Key,
			PartNumber: &partNumber,
			UploadId:   uploadID,
			Body:       bytes.NewReader(part),
		})
		if err != nil {
			abortInput := &s3.AbortMultipartUploadInput{
				Bucket:   &config.CONFIG.AWSBucket,
				Key:      &file.S3Key,
				UploadId: uploadID,
			}
			_, err := s3Client.AbortMultipartUpload(ctx, abortInput)
			if err != nil {
				return "", err
			}
			return "", err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadPartResp.ETag,
			PartNumber: &partNumber,
		})
		partNumber++
		if err == io.EOF {
			break
		}
	}

	output, err := s3Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   &config.CONFIG.AWSBucket,
		Key:      &file.S3Key,
		UploadId: uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return "", err
	}
	return strings.Trim(*output.ETag, strconv.Itoa(int('"'))), nil
}

func singlePartUpload(ctx context.Context, s3Client *s3.Client, file *models.File, data io.Reader) (string, error) {
	settings := config.CONFIG
	output, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &settings.AWSBucket,
		Key:           &file.S3Key,
		Body:          data,
		ContentType:   &file.ContentType,
		ContentLength: &file.Size,
	})
	if err != nil {
		return "", err
	}

	return strings.Trim(*output.ETag, strconv.Itoa(int('"'))), nil
}

func PutFile(ctx context.Context, file *models.File, data io.Reader) (string, error) {
	s3Client := InitS3()
	var eTag = ""
	var err error

	if file.Size <= MaxSingleUploadSize {
		eTag, err = singlePartUpload(ctx, s3Client, file, data)
	} else if file.Size >= MaxChunkSize {
		eTag, err = multiPartUpload(ctx, s3Client, file, data)
	}
	return eTag, err
}
