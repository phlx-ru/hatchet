package postgres

import (
	"context"
	"database/sql"
	"time"

	entDialectSQL "entgo.io/ent/dialect/sql"

	"github.com/phlx-ru/hatchet/metrics"
)

const (
	metricPrefixDefault = `postgres.connections`
)

type Database struct {
	driver         *entDialectSQL.Driver
	sendStatsEvery time.Duration
	metricPrefix   string
}

func (d *Database) Driver() *entDialectSQL.Driver {
	return d.driver
}

func (d *Database) DB() *sql.DB {
	return d.driver.DB()
}

func (d *Database) Check(ctx context.Context) error {
	_, err := d.DB().QueryContext(ctx, `select now()`)
	return err
}

func (d *Database) CollectDatabaseMetrics(ctx context.Context, metric metrics.Metrics) {
	if d.sendStatsEvery == 0 {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stats := d.DB().Stats()

		// The number of established connections both in use and idle.
		metric.Gauge(`postgres.connections.open`, stats.OpenConnections)

		// The number of connections currently in use.
		metric.Gauge(`postgres.connections.used`, stats.InUse)

		// The number of idle connections that is waiting for work.
		metric.Gauge(`postgres.connections.idle`, stats.Idle)

		// The total number of connections waited for.
		metric.Gauge(`postgres.connections.wait`, stats.WaitCount)

		// The total time blocked waiting for a new connection in milliseconds
		metric.Gauge(`postgres.connections.wait_duration`, stats.WaitDuration/time.Millisecond)

		// The total number of connections closed due to SetMaxIdleConns.
		metric.Gauge(`postgres.connections.max_idle_closed`, stats.MaxIdleClosed)

		// The total number of connections closed due to SetConnMaxIdleTime.
		metric.Gauge(`postgres.connections.max_idle_time_closed`, stats.MaxIdleTimeClosed)

		// The total number of connections closed due to SetConnMaxLifetime.
		metric.Gauge(`postgres.connections.max_lifetime_closed`, stats.MaxLifetimeClosed)

		time.Sleep(d.sendStatsEvery)
	}
}
