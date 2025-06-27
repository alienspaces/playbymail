package messenger

type Message struct {
	ID   string
	Body string
}

type SNSMessage struct {
	Message
	GroupID    string
	Subject    string
	Attributes map[string]string
}

// Messenger -
type Messenger interface {
	Publish(topic string, message SNSMessage) (messageID string, err error)
}
