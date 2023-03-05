package main

import (
	"context"
	"fmt"

	pubsub "github.com/eduardoths/pub-sub"
	"github.com/eduardoths/pub-sub/pkg/redis"
)

const TOPIC = "redis.test.general"

func main() {
	ctx := context.Background()
	app := pubsub.New(pubsub.Config{
		Listener: redis.NewRedisPubSub(ctx, redis.WithTopics(TOPIC)),
	})

	app.Route(TOPIC, firstMiddleware, secondMiddleware)
	if err := app.Listen(); err != nil {
		panic(err)
	}
}

func firstMiddleware(c *pubsub.Context) error {
	fmt.Println("First middleware")
	return c.Next()
}

func secondMiddleware(c *pubsub.Context) error {
	fmt.Printf("Second middleware. Message is: %v\n", string(c.Message.Data))

	return nil
}
