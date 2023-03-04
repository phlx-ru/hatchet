package watcher

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/alexcesaro/statsd.v2"
)

//go:generate moq -out dependencies_moq_test.go . logger baseLogger metrics

type logger interface {
	WithContext(ctx context.Context) *log.Helper
}

type baseLogger interface {
	log.Logger
}

type metrics interface {
	Increment(bucket string)
	NewTiming() statsd.Timing
}
