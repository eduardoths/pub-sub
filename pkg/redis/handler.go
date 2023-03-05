package redis

import (
	"context"
	"fmt"

	pubsub "github.com/eduardoths/pub-sub"
	"github.com/go-redis/redis/v8"
)

type RedisPubSub struct {
	client *redis.Client
	pubSub *redis.PubSub
	ctx    context.Context
}

func NewRedisPubSub(ctx context.Context, opts ...Option) *RedisPubSub {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Address, c.Port),
		Password: c.Password,
	})

	pubSub := redisClient.Subscribe(ctx, c.topics...)

	rps := &RedisPubSub{
		client: redisClient,
		pubSub: pubSub,
		ctx:    ctx,
	}
	return rps
}

func (rps *RedisPubSub) Listen(messages chan<- pubsub.Message, done <-chan error) error {
	for {
		select {
		case err := <-done:
			return err
		default:
			msg, err := rps.pubSub.ReceiveMessage(rps.ctx)
			if err != nil {
				return err
			}
			messages <- pubsub.Message{
				Topic: msg.Channel,
				Data:  []byte(msg.Payload),
			}
		}
	}
}
