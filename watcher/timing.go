package watcher

import "time"

type Timing interface {
	Send(bucket string)
	Duration() time.Duration
}

type EmptyTimingSendCall struct {
	Duration time.Duration
	Bucket   string
}

type EmptyTiming struct {
	start    time.Time
	bucket   string
	duration time.Duration
	history  []EmptyTimingSendCall
}

func NewEmptyTiming() *EmptyTiming {
	return &EmptyTiming{
		start: time.Now(),
	}
}

func (e *EmptyTiming) Send(bucket string) {
	e.bucket = bucket
	e.duration = e.Duration()
	e.history = append(e.history, EmptyTimingSendCall{
		Duration: e.Duration(),
		Bucket:   bucket,
	})
}

func (e *EmptyTiming) Duration() time.Duration {
	return time.Since(e.start)
}

func (e *EmptyTiming) History() []EmptyTimingSendCall {
	return e.history
}
