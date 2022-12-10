package pubsub

type Context struct {
	index    int
	handlers []Handler
	Message  Message
}

func (c *Context) Next() error {
	c.index++
	if len(c.handlers) > c.index {
		return c.handlers[c.index](c)
	}
	return nil
}
