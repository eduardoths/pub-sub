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

	t.Run("it should run multiple handlers", func(t *testing.T) {
		t.Parallel()
		const TOPIC_NAME = "example.topic.test"

		var (
			firstCounter  uint
			secondCounter uint
			firstMessage  = []byte("firstMessage")
			secondMessage = []byte("secondMessage")
		)

		firstHandler := func(c *pubsub.Context) error {
			firstCounter++
			assert.Equal(t, firstMessage, c.Message.Data)
			assert.Equal(t, TOPIC_NAME, c.Message.Topic)
			c.Message.Data = secondMessage
			return c.Next()
		}

		secondHandler := func(c *pubsub.Context) error {
			secondCounter++
			assert.Equal(t, secondMessage, c.Message.Data)
			assert.Equal(t, TOPIC_NAME, c.Message.Topic)
			return nil
		}

		mockListener := mocks.NewMockListener(pubsub.Message{
			Topic: TOPIC_NAME,
			Data:  firstMessage,
		})

		app := makeApp(mockListener)
		app.Route(TOPIC_NAME, firstHandler, secondHandler)
		assert.NoError(t, app.Listen())
		assert.EqualValues(t, 1, firstCounter)
		assert.EqualValues(t, 1, secondCounter)
	})

	t.Run("it should return nil if there are no more handlers", func(t *testing.T) {
		t.Parallel()
		const TOPIC_NAME = "example.topic.test"

		var (
			firstCounter uint
			firstMessage = []byte("firstMessage")
		)

		firstHandler := func(c *pubsub.Context) error {
			firstCounter++
			assert.Equal(t, firstMessage, c.Message.Data)
			assert.Equal(t, TOPIC_NAME, c.Message.Topic)
			return c.Next()
		}

		mockListener := mocks.NewMockListener(pubsub.Message{
			Topic: TOPIC_NAME,
			Data:  firstMessage,
		})

		app := makeApp(mockListener)
		app.Route(TOPIC_NAME, firstHandler)
		assert.NoError(t, app.Listen())
	})
}
