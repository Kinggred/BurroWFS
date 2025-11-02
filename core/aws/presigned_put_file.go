package aws

import (
	"burrowfs/core/config"
	"burrowfs/core/types"
	"context"
	"errors"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"time"
)

func GetPresignedURLs(files *[]types.FileDTO) (map[uuid.UUID]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s3Client := InitS3()
	if s3Client == nil {
		return nil, errors.New("s3 client not initialized")
	}
	s3PresignClient := s3.NewPresignClient(s3Client)
	var err error
	var presignedURLs map[uuid.UUID]string
	for _, file := range *files {
		var result *v4.PresignedHTTPRequest
		result, err = s3PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      &config.CONFIG.AWSBucket,
			Key:         &file.S3Key,
			ContentType: &file.ContentType,
		}, s3.WithPresignExpires(15*time.Minute))
		if err != nil {
			return nil, err
		}
		presignedURLs[*file.FileID] = result.URL
	}
	return presignedURLs, nil
}
