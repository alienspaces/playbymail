package sqsclient

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/aws/awsconfig"
	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

const (
	packageName = "sqsclient"
)

// NOTE:
// Overriding options (e.g., region) at operation call time is concurrency safe:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#OverrideClientOptionsForOperation
//
// The default retry mode is "standard", which means that the client will retry up to 3 times:
// https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
// Further confirmed by checking SQS client code.

// Client -
type Client struct {
	log logger.Logger
	sqs *sqs.Client

	config   Config
	messages []types.Message

	// receiveRequestAttemptID is used to receive the same messages during a visibility timeout,
	// if a networking issue occurs after SQS receives the ReceiveMessage API request,
	// but before the messages are acked.
	receiveRequestAttemptID *string
}

type NewFn func(context.Context, logger.Logger, Config) (messenger.SQS, error)

var _ messenger.SQS = &Client{}
var _ NewFn = New

// Config is required unless the field is a pointer type.
//
// The Region is the region of the SQS queue. If unspecified, the default AWS_REGION is used.
// Name or QueueURL must be specified.
type Config struct {
	Name     *string
	QueueURL *string
	Region   *string

	ConsumerConfig
	PublisherConfig
}

// ConsumerConfig may be specified if the SQS client is used for consuming messages.
type ConsumerConfig struct {
	VisibilityTimeoutSecs *int32
	WaitTimeSecs          *int32
}

// PublisherConfig may be specified if the SQS client is used for publishing messages.
type PublisherConfig struct{}

// New - The same client can be safely used to concurrently send multiple requests:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#concurrently-using-service-clients
//
// This package can be modified to accept and reuse the same SQS client across sqsclient.Client instances.
// However, there is no need at this stage.
func New(ctx context.Context, l logger.Logger, cfg Config) (messenger.SQS, error) {
	c := &Client{
		log: l,
	}
	l = c.logger("New")

	awscfg, err := awsconfig.Load(ctx, cfg.Region)
	if err != nil {
		l.Warn("failed to load AWS config >%v<", err)
		return nil, err
	}

	c.sqs = sqs.NewFromConfig(awscfg)

	c.config, err = c.resolveConfig(ctx, cfg)
	if err != nil {
		l.Warn("failed to resolve sqsclient config >%v<", err)
		return nil, err
	}

	return c, nil
}

func (c *Client) resolveConfig(ctx context.Context, cfg Config) (Config, error) {
	if cfg.Name == nil && cfg.QueueURL == nil {
		return Config{}, errors.New("either Name or QueueURL must be specified")
	} else if cfg.Name != nil && cfg.QueueURL == nil {
		url, err := c.getQueueURL(ctx, *cfg.Name)
		if err != nil {
			return Config{}, err
		}

		cfg.QueueURL = &url
	}

	cfg.VisibilityTimeoutSecs = resolveVisibilityTimeout(cfg.VisibilityTimeoutSecs)
	cfg.WaitTimeSecs = convert.Ptr(resolveWaitTime(cfg.WaitTimeSecs))

	return cfg, nil
}

func (c *Client) GetQueueURL() string {
	return convert.String(c.config.QueueURL)
}

// getQueueURL returns the URL of the queue for the specified name and region.
// name is the case-sensitive name of the queue, not the ARN or URL.
func (c *Client) getQueueURL(ctx context.Context, name string) (string, error) {
	l := c.logger("getQueueURL")

	urlResult, err := c.sqs.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &name,
	})
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		}

		return "", err
	}

	return convert.String(urlResult.QueueUrl), nil
}

func (c *Client) logger(funcName string) logger.Logger {
	if c.log == nil {
		return nil
	}

	l := c.log.WithPackageContext(packageName).WithFunctionContext(funcName)
	if c.config.QueueURL != nil {
		l.Context("queue-url", *c.config.QueueURL)
	}
	if c.config.Name != nil {
		l.Context("queue-name", *c.config.Name)
	}

	return l
}
