package pubsub

type Listener interface {
	Listen(messages chan<- Message, done <-chan error) error
}
