package message

import (
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

const (
	packageName = "message"
)

// Client -
type Client struct {
	log logger.Logger
	sns *sns.Client
}

var _ messenger.Messenger = &Client{}

// NewClient - The same client can be safely used to concurrently send multiple requests:
// https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#concurrently-using-service-clients
func NewClient(l logger.Logger) (*Client, error) {

	m := &Client{
		log: l,
	}

	err := m.Init()
	if err != nil {
		m.log.Warn("failed init >%v<", err)
		return nil, err
	}

	return m, nil
}

// Init -
func (m *Client) Init() error {
	l := m.logger("Init")

	snsClient, err := m.newSNSClient()
	if err != nil {
		l.Warn("Failed to create sns client >%v<", err)
		return err
	}
	m.sns = snsClient

	return nil
}

func (m *Client) logger(functionName string) logger.Logger {
	if m.log == nil {
		return nil
	}
	return m.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
