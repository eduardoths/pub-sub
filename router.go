package pubsub

type Router interface {
	Route(topic string, handlers ...Handler) Router
}

type Route struct {
	Handlers []Handler
}
