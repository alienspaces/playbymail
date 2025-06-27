package sqsclient

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

func (c *Client) Publish(ctx context.Context, msg messenger.SQSMessage) error {
	l := c.logger("Publish")

	_, err := c.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    &msg.Body,
		QueueUrl:       c.config.QueueURL,
		MessageGroupId: msg.GroupID,
	})
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			l.Warn("failed to publish message: >%v<", err)
		}

		return err
	}

	return nil
}
