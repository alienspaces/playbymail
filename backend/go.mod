module gitlab.com/alienspaces/playbymail

go 1.23.5

replace gitlab.com/alienspaces/playbymail/core => ./core

require (
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v1.29.17
	github.com/aws/aws-sdk-go-v2/service/s3 v1.81.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.7
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.7
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.8
	github.com/aws/aws-sdk-go-v2/service/ssm v1.59.3
	github.com/aws/smithy-go v1.22.4
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/caarlos0/env/v10 v10.0.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/klauspost/compress v1.18.0
	github.com/leekchan/accounting v1.0.0
	github.com/lib/pq v1.10.9
	github.com/r3labs/diff/v3 v3.0.1
	github.com/riverqueue/river v0.23.1
	github.com/riverqueue/river/riverdriver/riverpgxv5 v0.23.1
	github.com/rs/cors v1.11.1
	github.com/rs/zerolog v1.34.0
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	github.com/shopspring/decimal v1.4.0
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v2 v2.27.7
	github.com/xeipuuv/gojsonschema v1.2.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.11 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.70 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.34.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/riverqueue/river/riverdriver v0.23.1 // indirect
	github.com/riverqueue/river/rivershared v0.23.1 // indirect
	github.com/riverqueue/river/rivertype v0.23.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
