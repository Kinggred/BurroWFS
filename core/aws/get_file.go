package aws

import (
	"burrowfs/core/config"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetFileData(ctx context.Context, awsKey string) ([]byte, error) {
	s3Client := InitS3()
	if s3Client == nil {
		return nil, nil
	}
	return nil, nil
}

func GetFileURL(ctx context.Context, awsKey string) (string, error) {
	s3Client := InitS3()
	if s3Client == nil {
		return "", nil
	}

	s3PresignClient := s3.NewPresignClient(s3Client)
	presignResult, err := s3PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &config.CONFIG.AWSBucket,
		Key:    &awsKey,
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}
	return presignResult.URL, nil
}
