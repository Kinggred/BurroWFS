package aws

import (
	"burrowfs/core/config"
	"burrowfs/core/logging"
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3() *s3.Client {
	logger := logging.Get("AWS Initializer")
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
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
