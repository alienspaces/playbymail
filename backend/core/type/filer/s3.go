package filer

import "context"

type S3 interface {
	GetObject(ctx context.Context, key string) (data []byte, err error)
	PutObject(ctx context.Context, key string, data []byte) error
}
