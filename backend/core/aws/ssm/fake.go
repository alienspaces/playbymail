package ssmclient

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/secretor"
)

type FakeClient struct {
	secrets map[string]string
}

var _ secretor.SSM = &FakeClient{}
var _ NewFn = NewFake

func NewFake(_ context.Context, _ logger.Logger, _ Config) (secretor.SSM, error) {
	c := &FakeClient{
		secrets: map[string]string{},
	}
	return c, nil
}

func (c *FakeClient) GetSecret(_ context.Context, key string) (*string, error) {
	_, ok := c.secrets[key]
	if !ok {
		c.secrets[key] = record.NewRecordID()
	}

	return convert.Ptr(c.secrets[key]), nil
}
