package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/phlx-ru/hatchet/logger"
)

const (
	defaultEnv          = `local`
	defaultEnabled      = false
	defaultFlushTimeout = 2 * time.Second
)

var (
	defaultLogger = logger.NewHelper(log.DefaultLogger)
)

type Options struct {
	enabled      bool
	env          string
	flushTimeout time.Duration
	logger       logger.Logger
}

type Option func(o *Options)

func WithEnabled(enabled bool) Option {
	return func(o *Options) {
		o.enabled = enabled
	}
}

func WithEnv(env string) Option {
	return func(o *Options) {
		o.env = env
	}
}

func WithFlushTimeout(flushTimeout time.Duration) Option {
	return func(o *Options) {
		o.flushTimeout = flushTimeout
	}
}

func WithFlushTimeoutString(flushTimeoutDuration string) Option {
	return func(o *Options) {
		if duration, err := time.ParseDuration(flushTimeoutDuration); err == nil {
			o.flushTimeout = duration
		} else {
			o.logger.Errorf(`failed to parse sentry flush duration %s`, flushTimeoutDuration)
		}
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

func Init(dsn string, options ...Option) (func(), error) {
	if dsn == "" {
		return nil, nil
	}
	opts := &Options{
		enabled:      defaultEnabled,
		env:          defaultEnv,
		flushTimeout: defaultFlushTimeout,
		logger:       defaultLogger,
	}
	for _, option := range options {
		option(opts)
	}
	if !opts.enabled {
		return nil, nil
	}
	sentryOptions := sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      opts.env,
		TracesSampleRate: 1.0,
	}
	if err := sentry.Init(sentryOptions); err != nil {
		return nil, err
	}
	return flush(opts.flushTimeout, opts.logger), nil
}

func flush(duration time.Duration, logger logger.Logger) func() {
	return func() {
		sentry.Flush(duration)
		logger.Info(`sentry was flushed`)
	}
}
