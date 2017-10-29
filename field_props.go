package pqprobe

import "github.com/lib/pq/oid"

type (
	// Relation describes relation within database.
	Relation struct {
		Schema string `db:"schema"`
		Name   string `db:"name"`
		Type   string `db:"type"`
		Owner  string `db:"owner"`
	}

	// tableRelation contains table's oid.
	tableRelation struct {
		OID        int64  `db:"oid"`
		SchemaName string `db:"nspname"`
		TableName  string `db:"relname"`
	}

	// Field contains table's fields properties.
	Field struct {
		FieldNumber     int64  `db:"attnum"`
		TypeID          uint32 `db:"atttypid"`
		FieldName       string `db:"attname"`
		Type            string `db:"format_type"`
		Nullable        bool   `db:"attnotnull"`
		IsPrimary       bool   `db:"indisprimary"`
		IsUnique        bool   `db:"indisunique"`
		IsValid         bool   `db:"indisvalid"`
		IndexDefinition string `db:"pg_get_indexdef"`
	}
)

// Name returns field info in camelcase.
// Useful for struct field name.
func (t Field) Name() string {
	return camelify(variableNameRule(t.FieldName))
}

// GoTypeString returns this field's type as equivalent Golang type.
// Useful for struct field type.
// See https://godoc.org/github.com/lib/pq#hdr-Data_Types for conventions.
func (t Field) GoTypeString() string {
	return typeString(oid.Oid(t.TypeID))
}
