package sqsclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/collection/mmap"
)

func (c *Client) Ack(ctx context.Context) error {
	l := c.logger("Ack")
	cfg := c.config

	// For Standard SQS, there is no content-based deduplication. There may be duplicate messages in the same batch.
	// SQS batch delete requires that all message IDs are unique (otherwise a `BatchEntryIdsNotDistinct` error occurs).
	uniqueMessages := mmap.FromSlice(func(v types.Message) *string {
		return v.MessageId
	}, c.messages)

	var entries []types.DeleteMessageBatchRequestEntry
	for _, msg := range uniqueMessages {
		entries = append(entries, types.DeleteMessageBatchRequestEntry{
			Id:            msg.MessageId,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}

	if len(entries) == 0 {
		l.Error("no messages to delete: >%#v<", c) // Then how did we get here...?
		return nil
	}

	res, err := c.sqs.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		QueueUrl: cfg.QueueURL,
		Entries:  entries,
	})
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			err = fmt.Errorf("failed to call service: >%s<, operation: >%s<, error: >%w<", oe.Service(), oe.Operation(), oe.Unwrap())
		}

		l.Warn("failed to delete message batch: entries >%#v< messages >%#v< error >%v<", entries, c.messages, err)
		return err
	}

	// It is only possible to retry the same ReceiveMessage action if no messages have been deleted.
	if len(res.Successful) > 0 {
		c.clearReceiveRequestAttemptID()
	}

	if len(res.Failed) > 0 {
		return fmt.Errorf("failed to batch delete messages: >%#v<", res.Failed)
	}

	return nil
}
