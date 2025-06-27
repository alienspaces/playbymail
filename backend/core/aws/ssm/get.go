package ssmclient

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go"
)

func (c *Client) GetSecret(ctx context.Context, key string) (*string, error) {
	l := c.logger("GetSecret")

	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	res, err := c.client.GetParameter(ctx, input)
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			l.Warn("failed to get secret parameter >%#v< >%v<", input, err)
		}

		return nil, err
	}

	return res.Parameter.Value, nil
}
