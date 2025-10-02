package aws

import (
	"burrowfs/core/config"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetFileStream(ctx context.Context, awsKey string) (io.ReadCloser, error) {
	s3Client := InitS3()
	if s3Client == nil {
		return nil, nil
	}
	output, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &config.CONFIG.AWSBucket,
		Key:    &awsKey,
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
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
