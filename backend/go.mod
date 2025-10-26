module gitlab.com/alienspaces/playbymail

go 1.24.5

replace gitlab.com/alienspaces/playbymail/core => ./core

replace gitlab.com/alienspaces/playbymail/schema => ./schema

require (
	github.com/OpenPrinting/goipp v1.2.0
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/caarlos0/env/v10 v10.0.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/chromedp/cdproto v0.0.0-20250724212937-08a3db8b4327
	github.com/chromedp/chromedp v0.14.1
	github.com/davecgh/go-spew v1.1.1
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/klauspost/compress v1.18.0
	github.com/leekchan/accounting v1.0.0
	github.com/lib/pq v1.10.9
	github.com/otiai10/gosseract/v2 v2.4.1
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
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/go-json-experiment/json v0.0.0-20250725192818-e39067aee2d2 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
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
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
