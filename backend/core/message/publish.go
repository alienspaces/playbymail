package message

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"gitlab.com/alienspaces/playbymail/core/aws/awsconfig"
	"gitlab.com/alienspaces/playbymail/core/type/messenger"
)

// Publish -
func (m *Client) Publish(topicARN string, message messenger.SNSMessage) (messageID string, err error) {
	l := m.logger("Publish")

	awsMessageAttributes := map[string]types.MessageAttributeValue{}

	for k := range message.Attributes {
		v := message.Attributes[k]
		if v == "" {
			l.Debug("excluding message attribute k >%s< v >%s<, value is empty", k, v)
			continue
		}
		l.Debug("adding message attribute k >%s< v >%s<", k, v)
		awsMessageAttributes[k] = types.MessageAttributeValue{StringValue: &v, DataType: aws.String("String")}
	}

	l.Info("publishing message topic >%s< body >%s< attributes >%#v<", topicARN, message.Body, awsMessageAttributes)

	input := &sns.PublishInput{
		Message:           aws.String(message.Body),
		Subject:           aws.String(message.Subject),
		TopicArn:          aws.String(topicARN),
		MessageAttributes: awsMessageAttributes,
	}

	if strings.HasSuffix(topicARN, ".fifo") {
		input.MessageDeduplicationId = aws.String(message.ID)
		input.MessageGroupId = aws.String(message.GroupID)
	}

	result, err := m.sns.Publish(context.TODO(), input)
	if err != nil {
		l.Warn("failed publishing message >%v<", err)
		return "", err
	}

	// No message ID, failed send
	if *result.MessageId == "" {
		msg := "response does not contain message ID"
		l.Warn(msg)
		return "", fmt.Errorf(msg)
	}

	l.Info("publish result >%+v<", result)

	return *result.MessageId, nil
}

// newSNSClient
func (m *Client) newSNSClient() (*sns.Client, error) {

	cfg, err := awsconfig.Load(context.TODO(), nil)
	if err != nil {
		m.log.Warn("failed to load AWS config >%v<", err)
		return nil, err
	}

	client := sns.NewFromConfig(cfg)

	return client, nil
}
