package sqsclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

func (c *Client) Consume(ctx context.Context) ([]messenger.ConsumedSQSMessage, error) {
	l := c.logger("Consume")
	cfg := c.config

	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl: cfg.QueueURL,
		AttributeNames: []types.QueueAttributeName{
			"All",
			// For optimisation, can probably be reduced to the following:
			//"MessageGroupId",
			//"MessageDeduplicationId",
		},
		MaxNumberOfMessages: 10,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		ReceiveRequestAttemptId: c.resolveReceiveRequestAttemptID(), // used for FIFO, but ignored by standard SQS
		VisibilityTimeout:       *cfg.VisibilityTimeoutSecs,

		// This triggers long polling: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-short-and-long-polling.html#sqs-long-polling
		WaitTimeSeconds: *cfg.WaitTimeSecs, // WaitTimeSeconds must be less than HTTP client read timeout (default 30s)
	}

	res, err := c.sqs.ReceiveMessage(ctx, receiveMessageInput)
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			err = fmt.Errorf("failed to call service: >%s<, operation: >%s<, error: >%w<", oe.Service(), oe.Operation(), oe.Unwrap())
		}

		l.Warn("failed receiving messages: >%#v<", err)
		return nil, err
	}

	c.messages = res.Messages

	if len(res.Messages) == 0 {
		c.clearReceiveRequestAttemptID()
		l.Debug("no messages found")
		return nil, nil
	}

	var msgs []messenger.ConsumedSQSMessage
	for _, msg := range res.Messages {
		msgs = append(msgs, messenger.ConsumedSQSMessage{
			QueueURL:   convert.String(cfg.QueueURL),
			GroupID:    msg.Attributes["MessageGroupId"], // Only FIFO SQS messages may contain a `MessageGroupId`
			Attributes: msg.Attributes,
			Message: messenger.Message{
				ID:   convert.String(msg.MessageId),
				Body: convert.String(msg.Body),
			},
		})
	}

	return msgs, nil
}
