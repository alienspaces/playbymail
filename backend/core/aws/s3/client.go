package s3client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gitlab.com/alienspaces/playbymail/core/aws/awsconfig"
	"gitlab.com/alienspaces/playbymail/core/type/filer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

const (
	packageName = "s3client"
)

// NOTE:
// Overriding options (e.g., region) at operation call time is concurrency safe:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#OverrideClientOptionsForOperation
//
// The default retry mode is "standard", which means that the client will retry up to 3 times:
// https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html

// Client -
type Client struct {
	log logger.Logger
	s3  *s3.Client

	config Config
}

type NewFn func(context.Context, logger.Logger, Config) (filer.S3, error)

var _ filer.S3 = &Client{}
var _ NewFn = New

// Config is required unless the field is a pointer type.
//
// The Region is the region of the S3 Bucket. If unspecified, the default AWS_REGION is used.
// BucketName must be specified.
type Config struct {
	BucketName string  `env:"APP_S3_BUCKET_NAME"`
	Region     *string `env:"APP_S3_REGION"`
}

// New - The same client can be safely used concurrently:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#concurrently-using-service-clients
func New(ctx context.Context, l logger.Logger, cfg Config) (filer.S3, error) {
	c := &Client{
		log: l,
	}
	l = c.logger("New")

	awscfg, err := awsconfig.Load(ctx, cfg.Region)
	if err != nil {
		l.Warn("failed to load AWS config >%v<", err)
		return nil, err
	}

	c.s3 = s3.NewFromConfig(awscfg)

	c.config, err = c.resolveConfig(cfg)
	if err != nil {
		l.Warn("failed to resolve s3client config >%v<", err)
		return nil, err
	}

	return c, nil
}

func (c *Client) resolveConfig(cfg Config) (Config, error) {
	if cfg.BucketName == "" {
		return Config{}, errors.New("BucketName must be specified")
	}

	return cfg, nil
}

func (c *Client) logger(funcName string) logger.Logger {
	if c.log == nil {
		return nil
	}

	l := c.log.WithPackageContext(packageName).WithFunctionContext(funcName)
	l.Context("bucket-name", c.config.BucketName)

	return l
}
