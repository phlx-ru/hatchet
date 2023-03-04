package stubs

import "github.com/go-kratos/kratos/v2/log"

type Logger struct{}

func (l *Logger) Log(_ log.Level, _ ...interface{}) error {
	return nil
}

func NewLogger() log.Logger {
	return &Logger{}
}
