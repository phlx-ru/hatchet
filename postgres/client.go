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

func WithCreateIfNotExists(create bool) Option {
	return func(d *Database) {
		d.createIfNotExists = create
	}
}

func SendStatsEvery(duration time.Duration) Option {
	return func(d *Database) {
		d.sendStatsEvery = duration
	}
}

func Open(source string, options ...Option) (*Database, func(), error) {
	driver, err := entDialectSQL.Open(driverPostgres, source)
	if err != nil {
		return nil, nil, err
	}
	database := &Database{
		driver:       driver,
		metricPrefix: metricPrefixDefault,
	}
	for _, option := range options {
		option(database)
	}
	cleaning := cleanup(database.driver)
	if database.createIfNotExists {
		if err = createDatabaseIfNotExists(source); err != nil {
			return nil, cleaning, err
		}
	}
	return database, cleaning, nil
}

func cleanup(db *entDialectSQL.Driver) func() {
	return func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
}
