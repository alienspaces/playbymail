package s3client

import (
	"bytes"
	"context"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/compress"
	"gitlab.com/alienspaces/playbymail/core/type/filer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type FakeClient struct {
	log logger.Logger

	Data map[string][]byte
}

var _ filer.S3 = &FakeClient{}
var _ NewFn = NewFake

func NewFake(_ context.Context, l logger.Logger, _ Config) (filer.S3, error) {
	c := &FakeClient{
		log:  l,
		Data: map[string][]byte{},
	}
	return c, nil
}

func (c *FakeClient) PutObject(_ context.Context, key string, body []byte) error {
	l := c.log.WithFunctionContext("PutObject")

	c.Data[key] = compress.ZStdCompressBuffer(body, nil)

	l.Info("Stored key >%s< data len >%d<", key, len(c.Data[key]))

	return nil
}

func (c *FakeClient) GetObject(_ context.Context, key string) ([]byte, error) {
	l := c.log.WithFunctionContext("PutObject")

	l.Info("Getting key >%s< data len >%d<", key, len(c.Data[key]))

	// This fake implementation should align with the real Client.GetObject implementation
	r := strings.NewReader(string(c.Data[key]))

	var buf bytes.Buffer
	if err := compress.ZStdDecompressStream(r, &buf); err != nil {
		return nil, err
	}

	ret := buf.Bytes()
	l.Info("Returning len bytes >%d<", len(ret))

	return ret, nil
}
