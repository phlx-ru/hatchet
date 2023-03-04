package postgres

import (
	"database/sql"
	"fmt"
	"regexp"

	entDialectSQL "entgo.io/ent/dialect/sql"
)

const (
	databaseNameRegex      = `dbname=([a-zA-Z_][a-zA-Z0-9_]*)`
	defaultExistedDatabase = `postgres`
	databaseNamePostgres   = `dbname=` + defaultExistedDatabase
)

func createDatabaseIfNotExists(source string) error {
	original, err := extractDatabaseNameFromSource(source)
	if err != nil {
		return err
	}
	if original == defaultExistedDatabase {
		return nil // default database is always already created
	}
	db, closing, err := openDefaultDatabase(source)
	if err != nil {
		return err
	}
	defer closing()
	rows, err := db.Query(`select true as exists from pg_database where datname = $1`, original)
	if err != nil {
		return err
	}
	if rows.Next() {
		return nil
	}
	_, err = db.Exec(fmt.Sprintf(`create database %s`, original))
	return err
}

func extractDatabaseNameFromSource(source string) (string, error) {
	regex := regexp.MustCompile(databaseNameRegex)
	submatch := regex.FindAllStringSubmatch(source, 1)
	if len(submatch) == 0 {
		return "", fmt.Errorf(`data.database.source does not have submatch for dbname`)
	}
	match := submatch[0]
	if len(match) < 2 {
		return "", fmt.Errorf(`data.database.source does not have match for dbname`)
	}
	databaseName := match[1]
	if databaseName == "" {
		return "", fmt.Errorf(`data.database.source has empty dbname`)
	}
	return databaseName, nil
}

func openDefaultDatabase(source string) (*sql.DB, func(), error) {
	regex := regexp.MustCompile(databaseNameRegex)
	baseSource := regex.ReplaceAllString(source, databaseNamePostgres)
	db, err := entDialectSQL.Open(driverPostgres, baseSource)
	if err != nil {
		return nil, nil, err
	}
	return db.DB(), cleanup(db), nil
}
