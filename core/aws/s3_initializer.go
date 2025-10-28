package aws

import (
	"burrowfs/core/config"
	"burrowfs/core/logging"
	"context"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3() *s3.Client {
	logger := logging.Get("AWS Initializer")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if config.CONFIG.AWSAccessKeyID == "" || config.CONFIG.AWSSecretAccessKey == "" || config.CONFIG.AWSBucket == "" || config.CONFIG.AWSRegion == "" {
		logger.Fatal("AWS S3 is not configured properly, missing credentials or bucket/region")
		return nil
	}
	logger.Debug("Initializing AWS S3 client")
	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(config.CONFIG.AWSRegion),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.CONFIG.AWSAccessKeyID,
				config.CONFIG.AWSSecretAccessKey,
				"", // no session token
			),
		),
	)
	if err != nil {
		logger.Fatal("Unable to load AWS SDK config, " + err.Error())
	}

	return s3.NewFromConfig(cfg)
}
