package pubsub

type Handler = func(c *Context) error

type App struct {
	done     chan error
	stack    map[string]Route
	Listener Listener
}

func New() *App {
	return &App{
		stack: make(map[string]Route),
	}
}

func (a *App) Route(topic string, handlers ...Handler) Router {
	route, ok := a.stack[topic]
	if !ok {
		route = Route{Handlers: make([]Handler, 0, len(handlers))}
	}
	route.Handlers = append(route.Handlers, handlers...)
	a.stack[topic] = route
	return a
}

func (a *App) Listen() error {
	messages := make(chan Message)
	a.done = make(chan error, 1)
	go func() {
		a.Listener.Listen(messages, a.done)
	}()

	for {
		select {
		case err := <-a.done:
			return err
		}
	}
}

