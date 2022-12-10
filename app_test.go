package pubsub_test

import (
	"errors"
	"testing"

	pubsub "github.com/eduardoths/pub-sub"
	"github.com/eduardoths/pub-sub/mocks"
	"github.com/stretchr/testify/assert"
)

func makeApp(listener pubsub.Listener) *pubsub.App {
	return pubsub.New(pubsub.Config{
		Listener: listener,
	})
}

func TestApp_Listen(t *testing.T) {
	t.Parallel()
	t.Run("it should be shut down without returning error", func(t *testing.T) {
		t.Parallel()
		mockListener := mocks.NewMockListener().
			WithShutdown(nil)
		app := makeApp(mockListener)
		assert.NoError(t, app.Listen())
	})

	t.Run("it should be shut down returning an error", func(t *testing.T) {
		t.Parallel()
		err := errors.New("test-error")
		mockListener := mocks.NewMockListener().
			WithShutdown(err)
		app := makeApp(mockListener)
		assert.Equal(t, err, app.Listen())
	})

	t.Run("it should send messages to handler before shutting down", func(t *testing.T) {
		t.Parallel()
		const TOPIC_NAME = "example.topic"
		var callCount uint
		mockHandler := func(c *pubsub.Context) error {
			callCount += 1
			return nil
		}

		mockListener := mocks.NewMockListener(
			pubsub.Message{
				Topic: TOPIC_NAME,
			},
		)
		app := makeApp(mockListener)
		app.Route(TOPIC_NAME, mockHandler)
		assert.NoError(t, app.Listen())
		assert.Equal(t, uint(1), callCount, "handler called times mismatch")
	})

	t.Run("it should ignore messages if the topic isn't routered", func(t *testing.T) {
		t.Parallel()
		const TOPIC_NAME = "inexistent"

		mockListener := mocks.NewMockListener(
			pubsub.Message{
				Topic: TOPIC_NAME,
			},
		)
		app := makeApp(mockListener)
		assert.NoError(t, app.Listen())
	})
}
