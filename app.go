package pubsub

type Handler = func(c *Context) error

type App struct {
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
	done := make(chan error, 1)
	go func() {
		done <- a.Listener.Listen(messages, done)
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

	c := Context{}

	if len(route.Handlers) > 0 {
		route.Handlers[0](&c)
	}
}
