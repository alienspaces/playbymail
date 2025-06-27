package awsconfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Load automatically loads the following environment variables with the following precedence:
// 1. AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
// 2. AWS_WEB_IDENTITY_TOKEN_FILE
//
// This automatically loads the AWS_REGION environment variable, unless otherwise specified
// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
func Load(ctx context.Context, region *string, optFns ...func(o *config.LoadOptions) error) (aws.Config, error) {
	if region != nil {
		optFns = append(optFns, config.WithRegion(*region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
