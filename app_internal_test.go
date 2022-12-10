package pubsub

import (
	"testing"

	"github.com/eduardoths/pub-sub/tests/resources"
	"github.com/eduardoths/pub-sub/tests/resources/utils"
	"github.com/stretchr/testify/assert"
)

func makeApp() *App {
	return New(Config{
		Listener: nil,
	})
}

func TestApp_Route(t *testing.T) {
	t.Parallel()
	type in struct {
		app      *App
		topic    string
		handlers []Handler
	}

	type out struct {
		stack map[string]Route
	}

	fakeHandlers := []Handler{
		func(c *Context) error { return nil },
		func(c *Context) error { return nil },
	}

	ts := resources.UnitTestSuite[in, out]{
		{
			It: "should add a route without any handlers",
			In: in{
				app:      makeApp(),
				topic:    "xpto.topic.create",
				handlers: nil,
			},
			Want: out{
				stack: map[string]Route{
					"xpto.topic.create": {Handlers: make([]Handler, 0)},
				},
			},
		},
		{
			It: "should append one handler",
			In: in{
				app:      makeApp(),
				topic:    "xpto.topic.create",
				handlers: []Handler{fakeHandlers[0]},
			},
			Want: out{
				stack: map[string]Route{
					"xpto.topic.create": {Handlers: []Handler{fakeHandlers[0]}},
				},
			},
		},
		{
			It: "should append two handlers",
			In: in{
				app:      makeApp(),
				topic:    "xpto.topic.create",
				handlers: fakeHandlers[0:2],
			},
			Want: out{
				stack: map[string]Route{
					"xpto.topic.create": {Handlers: fakeHandlers[0:2]},
				},
			},
		},
		{
			It: "should append one handler at a route that already exists",
			In: in{
				app:      makeApp(),
				topic:    "xpto.topic.create",
				handlers: fakeHandlers[1:2],
			},
			Before: func(t *testing.T, in *in) {
				in.app.stack["xpto.topic.create"] = Route{
					Handlers: fakeHandlers[0:1],
				}
			},
			Want: out{
				stack: map[string]Route{
					"xpto.topic.create": {Handlers: fakeHandlers[0:2]},
				},
			},
		},
		{
			It: "should append handlers to another topic",
			In: in{
				app:      makeApp(),
				topic:    "xpto.topic.delete",
				handlers: fakeHandlers[1:2],
			},
			Before: func(t *testing.T, in *in) {
				in.app.stack["xpto.topic.create"] = Route{
					Handlers: fakeHandlers[0:1],
				}
				in.app.stack["xpto.topic.update"] = Route{
					Handlers: fakeHandlers[0:2],
				}

			},
			Want: out{
				stack: map[string]Route{
					"xpto.topic.create": {Handlers: fakeHandlers[0:1]},
					"xpto.topic.update": {Handlers: fakeHandlers[0:2]},
					"xpto.topic.delete": {Handlers: fakeHandlers[1:2]},
				},
			},
		},
	}

	for _, scenario := range ts {
		t.Run(scenario.It, func(t *testing.T) {
			if scenario.Before != nil {
				scenario.Before(t, &scenario.In)
			}
			actual := out{}
			scenario.In.app.Route(scenario.In.topic, scenario.In.handlers...)
			actual.stack = scenario.In.app.stack

			assert.Len(t, actual.stack, len(scenario.Want.stack))
			for k, v := range scenario.Want.stack {
				assertRoute(t, v, actual.stack[k])
			}
		})
	}
}

func assertRoute(t *testing.T, want Route, actual Route) bool {
	wantHandlers := make([]string, len(want.Handlers))
	actualHandlers := make([]string, len(actual.Handlers))

	for i := range wantHandlers {
		wantHandlers[i] = utils.FuncName(want.Handlers[i])
	}
	for i := range actualHandlers {
		actualHandlers[i] = utils.FuncName(actual.Handlers[i])
	}
	return assert.Equal(t, wantHandlers, actualHandlers)
}
