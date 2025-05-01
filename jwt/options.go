package jwt

import "time"

type Options struct {
	Period time.Duration
}

type Option func(*Options)

func WithPeriod(period time.Duration) Option {
	return func(options *Options) {
		options.Period = period
	}
}
