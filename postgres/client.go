package postgres

import (
	"time"

	entDialectSQL "entgo.io/ent/dialect/sql"
)

const (
	driverPostgres = `postgres`
)

type Option func(d *Database)

func WithMaxIdleConns(conns int) Option {
	return func(d *Database) {
		d.DB().SetMaxIdleConns(conns)
	}
}

func WithMaxOpenConns(conns int) Option {
	return func(d *Database) {
		d.DB().SetMaxOpenConns(conns)
	}
}

func WithConnMaxIdleTime(duration time.Duration) Option {
	return func(d *Database) {
		d.DB().SetConnMaxIdleTime(duration)
	}
}

func WithConnMaxLifetime(duration time.Duration) Option {
	return func(d *Database) {
		d.DB().SetConnMaxLifetime(duration)
	}
}

func WithMetricPrefix(metricPrefix string) Option {
	return func(d *Database) {
		d.metricPrefix = metricPrefix
	}
}

func SendStatsEvery(duration time.Duration) Option {
	return func(d *Database) {
		d.sendStatsEvery = duration
	}
}

func Open(source string, options ...Option) (database *Database, cleaning func(), err error) {
	driver, err := entDialectSQL.Open(driverPostgres, source)
	if err != nil {
		return nil, nil, err
	}
	database = &Database{
		driver:       driver,
		metricPrefix: metricPrefixDefault,
	}
	for _, option := range options {
		option(database)
	}
	return database, cleanup(database.driver), nil
}

func cleanup(db *entDialectSQL.Driver) func() {
	return func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
}

func CreateIfNotExists(source string, options ...DatabaseOption) error {
	return createDatabaseIfNotExists(source, options...)
}
