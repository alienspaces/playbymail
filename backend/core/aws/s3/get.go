package s3client

import (
	"bytes"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/compress"
)

func (c *Client) GetObject(ctx context.Context, key string) ([]byte, error) {
	l := c.logger("GetObject")

	res, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket:                  aws.String(c.config.BucketName),
		Key:                     aws.String(key),
		ResponseContentEncoding: aws.String("zstd"),
	})
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			l.Warn("failed to get object >%v<", err)
		}

		return nil, err
	}

	defer res.Body.Close()

	var buf bytes.Buffer
	if err := compress.ZStdDecompressStream(res.Body, &buf); err != nil {
		l.Warn("failed to decompress object body >%V<", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
