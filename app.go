package pubsub

type Handler func(c *Context) error

type App struct {
	stack    map[string]Route
	listener Listener
}

func New(config Config) *App {
	return &App{
		stack:    make(map[string]Route),
		listener: config.Listener,
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
	done := make(chan error, 1)
	go func() {
		done <- a.listener.Listen(messages, done)
	}()

	for {
		select {
		case err := <-done:
			return err
		case msg := <-messages:
			a.routeMessage(msg)
		}
	}
}

func (a *App) routeMessage(msg Message) {
	route, ok := a.stack[msg.Topic]
	if !ok {
		return
	}

	c := &Context{
		index:    -1,
		handlers: route.Handlers,
		Message:  msg,
	}
	c.Next()
}
