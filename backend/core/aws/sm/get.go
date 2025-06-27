package smclient

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
)

func (c *Client) GetSecret(ctx context.Context, input *secretsmanager.GetSecretValueInput) (*string, error) {
	l := c.logger("GetSecret")

	res, err := c.sm.GetSecretValue(ctx, input)
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			l.Warn("failed to get secret >%v<", err)
		}

		return nil, err
	}

	return res.SecretString, nil
}
