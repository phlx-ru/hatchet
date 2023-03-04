package watcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/require"
	"gopkg.in/alexcesaro/statsd.v2"
)

func mutedTiming() statsd.Timing {
	s, _ := statsd.New(statsd.Mute(true))
	return s.NewTiming()
}

func TestFluentInterfaceFunctionPurity(t *testing.T) {
	ctx := context.Background()
	instanceBasic := Make(`metric.prefix`)
	instanceGet := instanceBasic.OnPreparedMethod("Get").
		WithIgnoredErrors([]error{fmt.Errorf("error get")}).
		WithFields(map[string]any{
			"method": "get",
		})
	instanceGet.Results(func() (context.Context, error) {
		return ctx, fmt.Errorf("result error get")
	})
	instanceSet := instanceBasic.OnPreparedMethod("Set").
		WithIgnoredErrors([]error{fmt.Errorf("error set")}).
		WithFields(map[string]any{
			"method": "set",
		})
	instanceSet.Results(func() (context.Context, error) {
		return ctx, nil
	})
	require.Equal(t, "", instanceBasic.method)

	require.Equal(t, "get", instanceGet.method)
	require.Equal(t, "error get", instanceGet.ignoredErrors[0].Error())
	require.Equal(t, "get", instanceGet.fields["method"])

	require.Equal(t, "set", instanceSet.method)
	require.Equal(t, "error set", instanceSet.ignoredErrors[0].Error())
	require.Equal(t, "set", instanceSet.fields["method"])
}

func TestWatcher_Results_NoPanic(t *testing.T) {
	var err error
	defer Make(`metric.prefix`).OnMethod(`test`).WithTimings().Results(func() (context.Context, error) {
		return context.Background(), err
	})

	err = fmt.Errorf("some error")

	require.Error(t, err)
}

func TestWatcher_Results_Logging(t *testing.T) {
	metrics := &metricsMock{
		IncrementFunc: func(bucket string) {
			require.Equal(t, "metric.prefix.test.warning", bucket)
		},
		NewTimingFunc: mutedTiming,
	}

	baseLogger := &baseLoggerMock{
		LogFunc: func(level log.Level, keyvals ...interface{}) error {
			require.Equal(t, log.LevelWarn, level)
			require.Equal(t, 6, len(keyvals))
			require.Equal(t, "module", keyvals[0])
			require.Equal(t, "watcher_test", keyvals[1])
			require.Equal(t, log.DefaultMessageKey, keyvals[2])
			require.Equal(t, "metric.prefix has warning on test", keyvals[3])
			require.Equal(t, "custom_field1", keyvals[4])
			require.Equal(t, struct{ name string }{name: "custom value"}, keyvals[5])
			return nil
		},
	}

	logger := &loggerMock{
		WithContextFunc: func(ctx context.Context) *log.Helper {
			return log.NewHelper(log.With(baseLogger, `module`, `watcher_test`))
		},
	}

	warningMessage := `some warning`

	isWarning := func(err error) bool {
		return err.Error() == warningMessage
	}

	var err error
	defer Make(`metric.prefix`).
		OnPreparedMethod(`Test`).
		WithMetrics(metrics).
		WithLogger(logger).
		WithTimings().
		WithFields(map[string]any{
			"custom_field1": struct{ name string }{name: "custom value"},
		}).
		WithWarningErrorChecks([]func(error) bool{
			isWarning,
		}).
		Results(func() (context.Context, error) { return context.Background(), err })

	err = fmt.Errorf(warningMessage)

	require.Error(t, err)
}

type SleepingStruct struct {
	metrics metrics
	logger  logger
	watcher *Watcher
}

func NewSleepingStruct(metrics metrics, logger logger) *SleepingStruct {
	return &SleepingStruct{
		metrics: metrics,
		logger:  logger,
		watcher: New(`metric.prefix`, logger, metrics),
	}
}

func (s *SleepingStruct) SleepForAWhile(ctx context.Context, duration time.Duration, errorMessage string) error {
	var err error
	defer s.watcher.OnPreparedMethod(`SleepForAWhile`).WithTimings().WithFields(map[string]any{
		"duration": duration,
	}).Results(func() (context.Context, error) {
		return ctx, err
	})
	if errorMessage != "" {
		err = fmt.Errorf(errorMessage)
	}
	time.Sleep(duration)
	return err
}

func TestNew(t *testing.T) {
	ctx := context.Background()
	incrementCalls := 0
	metrics := &metricsMock{
		IncrementFunc: func(bucket string) {
			if incrementCalls == 0 {
				require.Equal(t, `metric.prefix.sleepForAWhile.failure`, bucket)
			} else {
				require.Equal(t, `metric.prefix.sleepForAWhile.success`, bucket)
			}
			incrementCalls++
		},
		NewTimingFunc: mutedTiming,
	}
	baseLogger := &baseLoggerMock{
		LogFunc: func(level log.Level, keyvals ...interface{}) error {
			m := mapFromKeyValues(keyvals)
			if incrementCalls == 0 {
				require.Equal(t, log.LevelError, level)
				require.Equal(t, 10, len(keyvals))
				require.Equal(t, 50*time.Millisecond, m["duration"])
			} else {
				require.Equal(t, log.LevelInfo, level)
				require.Equal(t, 6, len(keyvals))
				require.Equal(t, 10*time.Millisecond, m["duration"])
			}
			return nil
		},
	}
	logger := &loggerMock{
		WithContextFunc: func(ctx context.Context) *log.Helper {
			return log.NewHelper(log.With(baseLogger, `module`, `watcher_test_new`))
		},
	}
	sleeping := NewSleepingStruct(metrics, logger)
	err := sleeping.SleepForAWhile(ctx, 50*time.Millisecond, "test error message")
	require.Error(t, err)
	require.Equal(t, 1, len(metrics.NewTimingCalls()))
	require.Equal(t, 1, len(metrics.IncrementCalls()))
	require.Equal(t, 1, len(logger.WithContextCalls()))

	err = sleeping.SleepForAWhile(ctx, 10*time.Millisecond, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(metrics.NewTimingCalls()))
	require.Equal(t, 2, len(metrics.IncrementCalls()))
	require.Equal(t, 2, len(logger.WithContextCalls()))
}

func mapFromKeyValues(kvs []any) map[any]any {
	res := map[any]any{}
	for i := 0; i < len(kvs); i += 2 {
		res[kvs[i]] = kvs[i+1]
	}
	return res
}
