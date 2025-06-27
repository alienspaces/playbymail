package sqsclient

import (
	"github.com/google/uuid"

	"gitlab.com/alienspaces/playbymail/core/convert"
)

// For FIFO queues, ReceiveRequestAttemptId can be used to receive the same messages during a visibility timeout,
// if a networking issue occurs after SQS receives the ReceiveMessage API request, but before you receive the messages.
func (c *Client) resolveReceiveRequestAttemptID() *string {
	if c.receiveRequestAttemptID == nil {
		c.receiveRequestAttemptID = convert.Ptr(uuid.NewString())
		return c.receiveRequestAttemptID
	}
	return c.receiveRequestAttemptID
}

func (c *Client) clearReceiveRequestAttemptID() {
	c.receiveRequestAttemptID = nil
}

// AWS SQS default visibility timeout is 30s.
// The only time this is relevant is if the ACK fails multiple retries, or the container is killed before successfully ACKing.
var defaultVisibilityTimeoutSecs = convert.Ptr(int32(30))

const maxVisibilityTimeoutSecs = 12 * 60 * 60

func resolveVisibilityTimeout(timeoutSecs *int32) *int32 {
	if timeoutSecs == nil {
		return defaultVisibilityTimeoutSecs
	}

	timeout := *timeoutSecs
	if timeout < 0 {
		return defaultVisibilityTimeoutSecs
	} else if timeout > maxVisibilityTimeoutSecs {
		return defaultVisibilityTimeoutSecs
	}

	return timeoutSecs
}

const defaultWaitTimeSecs = 20
const maxWaitTimeSecs = 20

func resolveWaitTime(waitSecs *int32) int32 {
	if waitSecs == nil {
		return defaultWaitTimeSecs
	}

	timeout := *waitSecs
	if timeout < 0 {
		return defaultWaitTimeSecs
	} else if timeout > maxWaitTimeSecs {
		return defaultWaitTimeSecs
	}

	return timeout
}
