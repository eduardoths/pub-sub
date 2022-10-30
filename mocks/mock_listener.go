package mocks

import (
	pubsub "github.com/eduardoths/pub-sub"
)

type MockListener struct {
	done        chan error
	shutdownErr error
	messages    []pubsub.Message
}

func NewMockListener(messages ...pubsub.Message) *MockListener {
	return &MockListener{
		messages: messages,
	}
}

func (ml *MockListener) WithShutdown(err error) *MockListener {
	ml.shutdownErr = err
	return ml
}

func (ml *MockListener) Listen(messages chan<- pubsub.Message, done <-chan error) error {
	for {
		select {
		case err := <-done:
			return err
		default:
			for _, msg := range ml.messages {
				messages <- msg
			}
			return ml.shutdownErr
		}
	}
}
