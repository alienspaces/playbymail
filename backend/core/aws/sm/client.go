package smclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"gitlab.com/alienspaces/playbymail/core/aws/awsconfig"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/secretor"
)

const (
	packageName = "smclient"
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
	sm  *secretsmanager.Client

	config Config
}

type NewFn func(context.Context, logger.Logger, Config) (secretor.SM, error)

var _ secretor.SM = &Client{}
var _ NewFn = New

// Config is required unless the field is a pointer type.
//
// The Region is the region of the Secrets Manager. If unspecified, the default AWS_REGION is used.
type Config struct {
	Region *string
}

// New - The same client can be safely used concurrently:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#concurrently-using-service-clients
func New(ctx context.Context, l logger.Logger, cfg Config) (secretor.SM, error) {
	c := &Client{
		log: l,
	}
	l = c.logger("New")

	awscfg, err := awsconfig.Load(ctx, cfg.Region)
	if err != nil {
		l.Warn("failed to load AWS config >%v<", err)
		return nil, err
	}

	c.sm = secretsmanager.NewFromConfig(awscfg)

	return c, nil
}

func (c *Client) logger(funcName string) logger.Logger {
	if c.log == nil {
		return nil
	}

	l := c.log.WithPackageContext(packageName).WithFunctionContext(funcName)

	return l
}
