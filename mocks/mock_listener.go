package mocks

import pubsub "github.com/eduardoths/pub-sub"

type MockListener struct {
	done      chan error
	doneError error
	messages  []pubsub.Message
}

func NewMockListener(messages ...pubsub.Message) *MockListener {
	return &MockListener{
		messages: messages,
	}
}

func (ml *MockListener) WithShutdown(err error) *MockListener {
	ml.doneError = err
	return ml
}

func (ml *MockListener) Listen(messages chan<- pubsub.Message, done chan error) {
	ml.done = done
	for {
		select {
		case <-ml.done:
			return
		default:
			for _, msg := range ml.messages {
				messages <- msg
			}
			ml.done <- ml.doneError
		}
	}
}

func (ml *MockListener) Shutdown(err error) {
	ml.doneError = err
}
