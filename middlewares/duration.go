package middlewares

import (
	"context"
	"strings"

	"github.com/phlx-ru/hatchet/metrics"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Duration(metric metrics.Metrics, prefix string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			slug := metricSlug(ctx)
			timing := metric.NewTiming()
			defer func() {
				if slug != "" {
					prefix = strings.TrimSuffix(prefix, `.`)
					timing.Send(prefix + `.` + slug)
				}
			}()
			return handler(ctx, req)
		}
	}
}

func metricSlug(ctx context.Context) string {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ""
	}
	slug := strings.ToLower(tr.Operation())
	slug = strings.ReplaceAll(strings.ToLower(slug), ".", "_")
	slug = strings.ReplaceAll(slug, "/", ".")
	kind := tr.Kind().String()
	return kind + slug
}
