package messenger

import "context"

// SQSMessage -
//
// GroupID should only be populated for messages to FIFO SQS.
type SQSMessage struct {
	Body    string
	GroupID *string
}

// ConsumedSQSMessage is used for consumed SQS messages.
//
// GroupID may be an empty string as the SQS message may not contain a `MessageGroupId`.
type ConsumedSQSMessage struct {
	Message
	QueueURL   string
	GroupID    string
	Attributes map[string]string
}

type SQS interface {
	Consume(ctx context.Context) (message []ConsumedSQSMessage, err error)
	Ack(ctx context.Context) (err error)
	Publish(ctx context.Context, msg SQSMessage) error
	GetQueueURL() string
}
