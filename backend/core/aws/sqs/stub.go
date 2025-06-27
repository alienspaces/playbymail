package sqsclient

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

type StubClient struct {
	URL          string
	Data         []messenger.ConsumedSQSMessage
	WaitTimeSecs *int32
}

var _ messenger.SQS = &StubClient{}
var _ NewFn = NewStub

func NewStub(_ context.Context, l logger.Logger, cfg Config) (messenger.SQS, error) {
	if cfg.QueueURL == nil {
		return nil, fmt.Errorf("queue URL is required")
	}

	c := &StubClient{
		URL:          *cfg.QueueURL,
		WaitTimeSecs: cfg.WaitTimeSecs,
	}

	return c, nil
}

func (c *StubClient) GetQueueURL() string {
	return c.URL
}

func (c *StubClient) Consume(ctx context.Context) ([]messenger.ConsumedSQSMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if c.WaitTimeSecs != nil {
				time.Sleep(time.Duration(*c.WaitTimeSecs) * time.Second)
			}

			d := c.Data
			c.Data = nil
			return d, nil
		}
	}
}

func (c *StubClient) Publish(ctx context.Context, msg messenger.SQSMessage) error {
	return nil
}

func (c *StubClient) Ack(ctx context.Context) error {
	c.Data = nil
	return nil
}
