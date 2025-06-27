package smclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/secretor"
)

type FakeClient struct {
	secrets map[string]string
}

var _ secretor.SM = &FakeClient{}
var _ NewFn = NewFake

func NewFake(_ context.Context, _ logger.Logger, _ Config) (secretor.SM, error) {
	c := &FakeClient{
		secrets: map[string]string{},
	}
	return c, nil
}

func (c *FakeClient) GetSecret(_ context.Context, input *secretsmanager.GetSecretValueInput) (*string, error) {
	_, ok := c.secrets[*input.SecretId]
	if !ok {
		c.secrets[*input.SecretId] = record.NewRecordID()
	}

	return convert.Ptr(c.secrets[*input.SecretId]), nil
}
