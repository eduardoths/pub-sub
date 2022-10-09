package pubsub

type Handler = func(c *Context) error

type App struct {
	stack map[string]Route
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
