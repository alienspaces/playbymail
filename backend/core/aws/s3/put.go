package s3client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"gitlab.com/alienspaces/playbymail/core/compress"
)

// Since ServiceGate settle requests will only be 8MB, there is no need to use multipart upload.
// AWS provides an SDK that implements multipart upload (which keeps track of part numbers and ETags).
// AWS recommends using multipart upload for files larger than 100MB.
// By default, on network failure, uploaded parts are deleted.
// To enable resuming multipart uploads, we need to disable this behaviour, and keep track of the S3 Object ID ourselves.

func (c *Client) PutObject(ctx context.Context, key string, body []byte) error {
	l := c.logger("Upload")

	compressed := compress.ZStdCompressBuffer(body, nil)

	// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
	hash := md5.Sum(compressed)
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:          aws.String(c.config.BucketName),
		Key:             aws.String(key),
		Body:            bytes.NewReader(compressed),
		ContentMD5:      aws.String(base64.StdEncoding.EncodeToString(hash[:])), // Needed to upload an object with a retention period configured using Amazon S3 Object Lock. Also recommended as an e2e integrity check.
		ContentEncoding: aws.String("zstd"),
	})
	if err != nil {
		var oe *smithy.OperationError // wraps context cancelled error
		if errors.As(err, &oe) {
			l.Warn("failed to call service: >%s<, operation: >%s<, error: >%v<", oe.Service(), oe.Operation(), oe.Unwrap())
		} else {
			l.Warn("failed to put object >%v<", err)
		}

		return err
	}

	return nil
}
