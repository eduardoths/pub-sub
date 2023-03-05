package redis

import (
	"log"

	"github.com/caarlos0/env"
)

type config struct {
	Address  string `env:"REDIS_ADDR" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PSWD"`
	topics   []string
}

func newConfig() *config {
	rc := new(config)
	if err := env.Parse(rc); err != nil {
		log.Panicf("An unexpected error occurred. %s", err.Error())
	}
	return rc
}

type Option func(*config)

func WithTopics(topics ...string) Option {
	return func(c *config) {
		c.topics = topics
	}
}
