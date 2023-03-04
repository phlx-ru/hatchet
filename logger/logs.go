package logger

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

const (
	messageKey = "msg"
)

type Logger interface {
	WithContext(ctx context.Context) *log.Helper
	Log(level log.Level, keyvals ...interface{})
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Debugw(keyvals ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Infow(keyvals ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Warnw(keyvals ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Errorw(keyvals ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalw(keyvals ...interface{})
}

func SetGlobalDefaultLogger(level string) {
	loggerInstance := log.With(
		log.DefaultLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logLevel := log.ParseLevel(level)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(logLevel))
	log.SetLogger(logger)
}

type Log struct {
	logger      log.Logger
	sentryLevel log.Level
}

func (l *Log) Log(level log.Level, keyvals ...interface{}) error {
	err := l.logger.Log(level, keyvals...)
	if level >= l.sentryLevel {
		tags := ExtractMapFromKeyvals(keyvals...)
		if msg, ok := tags[messageKey]; ok {
			delete(tags, messageKey)
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetExtras(tags)
				message := fmt.Sprintf("%s: %s", level.String(), msg)
				if level == log.LevelError || level == log.LevelFatal {
					sentry.CaptureException(fmt.Errorf(message))
				} else {
					sentry.CaptureMessage(message)
				}
			})
		}
	}
	return err
}

func ExtractMapFromKeyvals(keyvals ...interface{}) map[string]any {
	res := map[string]any{}
	if len(keyvals) == 0 {
		return res
	}
	if len(keyvals) == 1 {
		res[messageKey] = keyvals[0]
		return res
	}
	length := len(keyvals)
	for i := 0; i < length; i += 2 {
		key := fmt.Sprintf("%v", keyvals[i])
		var value any
		if i != length-1 {
			value = keyvals[i+1]
		}
		res[key] = value
	}
	return res
}

func NewWithSentry(level string) log.Logger {
	return &Log{
		logger:      log.DefaultLogger,
		sentryLevel: log.ParseLevel(level),
	}
}

func New(id, name, version, env, level, sentryLevel string) *log.Filter {
	loggerWithSentry := NewWithSentry(sentryLevel)
	loggerInstance := log.With(
		loggerWithSentry,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"env", env,
		"service.id", id,
		"service.name", name,
		"service.version", version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logLevel := log.ParseLevel(level)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(logLevel))
	log.SetLogger(logger)
	return logger
}

func NewHelper(logger log.Logger, kv ...interface{}) *log.Helper {
	return log.NewHelper(log.With(logger, kv...))
}
