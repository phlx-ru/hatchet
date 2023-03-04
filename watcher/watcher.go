package watcher

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	pkgStrings "github.com/phlx-ru/hatchet/strings"
)

const (
	keyError      = `error`
	keyStack      = `stack`
	chunkSuccess  = `success`
	chunkWarning  = `warning`
	chunkFailure  = `failure`
	chunkTimings  = `timings`
	methodUnknown = `unknown`
)

type Watcher struct {
	logger              logger
	metrics             metrics
	method              string
	metricPrefix        string
	timing              Timing
	fields              map[string]any
	ignoredErrors       []error
	ignoredErrorsChecks []func(error) bool
	warningErrors       []error
	warningErrorsChecks []func(error) bool
}

// New builds watcher, that can use for a struct variable
// Usage:
//
//	type SleepingStruct struct {
//		metrics metrics
//		logger  logger
//		watcher *Watcher
//	}
//
//	func NewSleepingStruct(metrics metrics, logger logger) *SleepingStruct {
//		return &SleepingStruct{
//			metrics: metrics,
//			logger:  logger,
//			watcher: New(`metric.prefix`, logger, metrics),
//		}
//	}
//
//	func (s *SleepingStruct) SleepForAWhile(ctx context.Context, duration time.Duration, errorMessage string) error {
//		var err error
//		defer s.watcher.OnPreparedMethod(`SleepForAWhile`).WithTimings().Results(func() (context.Context, error) {
//			return ctx, err
//		})
//		if errorMessage != "" {
//			err = fmt.Errorf(errorMessage)
//		}
//		time.Sleep(duration)
//		return err
//	}
func New(metricPrefix string, logger logger, metrics metrics) *Watcher {
	return &Watcher{
		logger:       logger,
		metrics:      metrics,
		metricPrefix: metricPrefix,
	}
}

// Make makes fluent interface base with only metric prefix.
// Usage:
//
//	func (s *SleepingStruct) SleepForAWhile(ctx context.Context, duration time.Duration, errorMessage string) error {
//		var err error
//		defer Make(`services.prefix`).
//			OnPreparedMethod(`SleepForAWhile`).
//			WithLogger(s.logger).
//			WithMetrics(s.metrics).
//			WithTimings().
//			Results(func() (context.Context, error) { return ctx, err })
//		if errorMessage != "" {
//			err = fmt.Errorf(errorMessage)
//		}
//		time.Sleep(duration)
//		return err
//	}
func Make(metricPrefix string) *Watcher {
	return &Watcher{
		metricPrefix: metricPrefix,
	}
}

func (w *Watcher) WithMetrics(metrics metrics) *Watcher {
	n := *w
	n.metrics = metrics
	return &n
}

func (w *Watcher) WithLogger(logger logger) *Watcher {
	n := *w
	n.logger = logger
	return &n
}

// WithTimings add timings. WARNING: if call it before WithMetrics() then empty timing will be used
func (w *Watcher) WithTimings() *Watcher {
	n := *w
	if n.metrics != nil {
		n.timing = n.metrics.NewTiming()
	} else {
		n.timing = NewEmptyTiming()
	}
	return &n
}

func (w *Watcher) WithFields(fields map[string]any) *Watcher {
	n := *w
	n.fields = fields
	return &n
}

func (w *Watcher) WithIgnoredErrors(ignoredErrors []error) *Watcher {
	n := *w
	n.ignoredErrors = ignoredErrors
	return &n
}

func (w *Watcher) WithIgnoredErrorsChecks(ignoredErrorsChecks []func(error) bool) *Watcher {
	n := *w
	n.ignoredErrorsChecks = ignoredErrorsChecks
	return &n
}

func (w *Watcher) WithWarningErrors(warningErrors []error) *Watcher {
	n := *w
	n.warningErrors = warningErrors
	return &n
}

func (w *Watcher) WithWarningErrorChecks(warningErrorsChecks []func(error) bool) *Watcher {
	n := *w
	n.warningErrorsChecks = warningErrorsChecks
	return &n
}

func prepareMethodForMetric(method string) string {
	if method == "" {
		return ""
	}
	return strings.ToLower(method[:1]) + method[1:]
}

func (w *Watcher) OnPreparedMethod(method string) *Watcher {
	n := *w
	n.method = prepareMethodForMetric(method)
	return &n
}

func (w *Watcher) OnMethod(method string) *Watcher {
	n := *w
	n.method = method
	return &n
}

type ContextAndErrorCatcher func() (context.Context, error)

func (w *Watcher) Results(catcher ContextAndErrorCatcher) {
	ctx, err := catcher()
	result := chunkSuccess
	isWarning := isWarningError(err, w.warningErrors, w.warningErrorsChecks)
	isIgnored := isIgnoredError(err, w.ignoredErrors, w.ignoredErrorsChecks)
	if err != nil {
		result = chunkFailure
		if isWarning {
			result = chunkWarning
		}
		if isIgnored {
			result = chunkSuccess
		}
	}
	if w.logger != nil {
		if w.method == "" {
			w.logger.WithContext(ctx).Errorf("empty 'method' on watcher for metric prefix '%s'", w.metricPrefix)
			w.method = methodUnknown
		}
		if w.fields == nil {
			w.fields = map[string]any{}
		}
		w.fields[log.DefaultMessageKey] = fmt.Sprintf(`%s has %s on %s`, w.metricPrefix, result, w.method)
		if err != nil && !isIgnored && !isWarning {
			w.fields[keyError] = err
			w.fields[keyStack] = string(debug.Stack())
		}
		kvs := sortedKeyValuesFromFields(w.fields)
		w.fields = nil
		if err != nil && !isIgnored {
			if isWarning {
				w.logger.WithContext(ctx).Warnw(kvs...)
			} else {
				w.logger.WithContext(ctx).Errorw(kvs...)
			}
		} else {
			w.logger.WithContext(ctx).Infow(kvs...)
		}
	}
	metricStarts := w.method
	if w.metricPrefix != "" {
		metricStarts = pkgStrings.Metric(w.metricPrefix, w.method)
	}
	if w.timing != nil {
		w.timing.Send(pkgStrings.Metric(metricStarts, chunkTimings, result))
	}
	if w.metrics != nil {
		w.metrics.Increment(pkgStrings.Metric(metricStarts, result))
	}
}

func isErrorMatched(err error, matches []error, checks []func(error) bool) bool {
	if err == nil {
		return false
	}
	for _, matched := range matches {
		if errors.Is(err, matched) {
			return true
		}
	}
	for _, check := range checks {
		if check(err) {
			return true
		}
	}
	return false
}

func isIgnoredError(err error, ignoredErrors []error, ignoredErrorsChecks []func(error) bool) bool {
	return isErrorMatched(err, ignoredErrors, ignoredErrorsChecks)
}

func isWarningError(err error, warningErrors []error, warningErrorsChecks []func(error) bool) bool {
	return isErrorMatched(err, warningErrors, warningErrorsChecks)
}

func sortedKeyValuesFromFields(fields map[string]any) []any {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == log.DefaultMessageKey {
			return true
		}
		if keys[j] == log.DefaultMessageKey {
			return false
		}
		if keys[i] == keyStack {
			return false
		}
		if keys[j] == keyStack {
			return true
		}
		if keys[i] == keyError {
			return false
		}
		if keys[j] == keyError {
			return true
		}
		return false
	})
	kvs := make([]any, 0, len(fields)*2)
	for _, key := range keys {
		kvs = append(kvs, key, fields[key])
	}
	return kvs
}
