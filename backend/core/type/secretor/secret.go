package secretor

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SM interface {
	GetSecret(ctx context.Context, input *secretsmanager.GetSecretValueInput) (secret *string, err error)
}

type SSM interface {
	GetSecret(ctx context.Context, key string) (secret *string, err error)
}
