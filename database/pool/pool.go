// Package pool implements high level database pool for application.
// Application can store and retrieve database connections through this package, eliminating the
// needs to store database connection directly in package that needs them.
package pool

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	dbpool map[string]*sqlx.DB

	// ErrLabelExists is returned when label clash occurs.
	ErrLabelExists = errors.New("pool: label already exists")
	// ErrNilDB is returned when supplied database is nil.
	ErrNilDB = errors.New("pool: nil db")
)

// Open wraps sqlx.Open by also storing new connection into the pool under given label.
// Returns ErrLabelExists if another database connection is already stored under same label.
func Open(driverName, dataSourceName, label string) (db *sqlx.DB, err error) {
	db, err = sqlx.Open(driverName, dataSourceName)
	return db, Add(db, label)
}

// Add stores database connection under given label.
// Returns ErrLabelExists if another database connection is already stored under same label.
// To overwrite database connection with same label, use ForceAdd.
func Add(db *sqlx.DB, label string) (err error) {
	if db == nil {
		return ErrNilDB
	}
	if _, ok := dbpool[label]; !ok {
		dbpool[label] = db
		return nil
	}
	return ErrLabelExists
}

// ForceAdd stores database connection under given label.
// ForceAdd doesn't check for label clash and will overwrite connection under same label.
func ForceAdd(db *sqlx.DB, label string) (err error) {
	if db == nil {
		return ErrNilDB
	}
	dbpool[label] = db
	return nil
}

// Get returns database connection with given label.
func Get(label string) (db *sqlx.DB) {
	if db, ok := dbpool[label]; ok {
		return db
	}
	return nil
}
