package pubsub_test

import (
	"errors"
	"testing"

	pubsub "github.com/eduardoths/pub-sub"
	"github.com/eduardoths/pub-sub/mocks"
	"github.com/stretchr/testify/assert"
)

func TestApp_Listen(t *testing.T) {
	t.Parallel()
	t.Run("it should be shut down without returning error", func(t *testing.T) {
		t.Parallel()
		app := pubsub.New()
		mockListener := mocks.NewMockListener().
			WithShutdown(nil)
		app.Listener = mockListener
		assert.NoError(t, app.Listen())
	})

	t.Run("it should be shut down returning an error", func(t *testing.T) {
		t.Parallel()
		app := pubsub.New()
		err := errors.New("test-error")
		mockListener := mocks.NewMockListener().
			WithShutdown(err)
		app.Listener = mockListener
		assert.Equal(t, err, app.Listen())
	})

}
