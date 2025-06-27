package emailer

// Emailer -
type Emailer interface {
	Send(*Message) error
}
