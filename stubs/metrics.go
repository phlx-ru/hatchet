package stubs

import (
	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/phlx-ru/hatchet/metrics"
)

func NewMetrics() metrics.Metrics {
	s, err := statsd.New(statsd.Mute(true))
	if err != nil {
		panic(err)
	}
	return s
}
