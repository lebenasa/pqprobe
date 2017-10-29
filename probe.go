package pqprobe

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	// Prober is an interface to table & fields discovery functions.
	Prober interface {
		QueryRelations() (relations []Relation, err error)
		QueryFields(tableName string) (fields []Field, err error)
	}

	// pqProber enables table & fields discovery for postgresql database.
	pqProber struct {
		db                  *sqlx.DB
		selectRelations     *sqlx.Stmt
		selectTableRelation *sqlx.Stmt
		selectFieldProps    *sqlx.Stmt
	}
)

// Open wraps sqlx.Open to return a Prober.
// Will return error if Prober implementation for given driver is not yet implemented.
// Currently supported driver:
// 	- postgres
func Open(driverName, dataSourceName string) (prober Prober, err error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open connection to db")
	}

	switch driverName {
	case "postgres":
		return NewPqProber(db)
	}

	return nil, errors.New("unsupported driver")
}

// NewPqProber wraps given postgresql database into Prober to discover its table & fields information.
func NewPqProber(db *sqlx.DB) (prober Prober, err error) {
	if db.DriverName() != "postgres" {
		return nil, errors.New("mismatch sql driver")
	}

	selectDatabaseRelation, err := db.Preparex(`
		SELECT n.nspname as "Schema",
		  c.relname as "Name",
		  CASE c.relkind WHEN 'r' THEN 'table' WHEN 'v' THEN 'view' WHEN 'm' 
		  	THEN 'materialized view' WHEN 'i' THEN 'index' WHEN 'S' THEN 'sequence' WHEN 's' 
			THEN 'special' WHEN 'f' THEN 'foreign table' END as "Type",
		  pg_catalog.pg_get_userbyid(c.relowner) as "Owner"
		FROM pg_catalog.pg_class c
		     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind IN ('r','v','m','S','f','')
		      AND n.nspname <> 'pg_catalog'
		      AND n.nspname <> 'information_schema'
		      AND n.nspname !~ '^pg_toast'
		  AND pg_catalog.pg_table_is_visible(c.oid)
		ORDER BY 1,2;
	`)
	if err != nil {
		return nil, errors.Wrap(err, "prepare failed")
	}

	selectTableRelation, err := db.Preparex(`
		SELECT 
		  c.oid,
		  n.nspname,
		  c.relname
		FROM pg_catalog.pg_class c
		     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relname ~ $1
		  AND pg_catalog.pg_table_is_visible(c.oid)
		ORDER BY 2, 3;
	`)
	if err != nil {
		return nil, errors.Wrap(err, "prepare failed")
	}

	selectFieldProps, err := db.Preparex(`
		SELECT
		  a.attnum,
		  a.atttypid,
		  a.attname,
		  pg_catalog.format_type(a.atttypid, a.atttypmod),
		  a.attnotnull,
		  COALESCE(i.indisprimary, false) indisprimary,
		  COALESCE(i.indisunique, false) indisunique,
		  COALESCE(i.indisvalid, false) indisvalid,
		  COALESCE(pg_catalog.pg_get_indexdef(i.indexrelid, 0, true), '') pg_get_indexdef
		FROM pg_catalog.pg_attribute a
		  LEFT JOIN pg_catalog.pg_index i ON (i.indrelid = a.attrelid AND a.attnum = ANY (i.indkey))
		WHERE a.attrelid = $1 AND a.attnum > 0 AND NOT a.attisdropped
		ORDER BY a.attnum;
	`)
	if err != nil {
		return nil, errors.Wrap(err, "prepare failed")
	}

	return pqProber{db, selectDatabaseRelation, selectTableRelation, selectFieldProps}, nil
}

// QueryRelations probes the database for all relations within it.
func (p pqProber) QueryRelations() (relations []Relation, err error) {
	relationRows, err := p.selectRelations.Queryx()
	if err != nil {
		return nil, errors.Wrap(err, "relation probe failed")
	}
	defer relationRows.Close()

	for relationRows.Next() {
		r := Relation{}
		err = relationRows.StructScan(&r)
		if err != nil {
			return nil, errors.Wrap(err, "struct scan failed on relation probe")
		}
		relations = append(relations, r)
	}

	return
}

// QueryFields probes the database for given table name and returns its fields' properties.
func (p pqProber) QueryFields(tableName string) (fields []Field, err error) {
	rel := tableRelation{}
	err = p.selectTableRelation.QueryRowx(fmt.Sprintf("^(%v)$", tableName)).StructScan(&rel)
	if err != nil {
		return nil, errors.Wrapf(err, "table %v probe failed", tableName)
	}

	fieldRows, err := p.selectFieldProps.Queryx(rel.OID)
	if err != nil {
		return nil, errors.Wrapf(err, "fields probe failed for table %v", tableName)
	}
	defer fieldRows.Close()

	for fieldRows.Next() {
		ti := Field{}
		err = fieldRows.StructScan(&ti)
		if err != nil {
			return nil, errors.Wrapf(err, "struct scan failed on probe for table %v", tableName)
		}
		fields = append(fields, ti)
	}

	return
}
