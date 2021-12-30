package metadata

import (
	"github.com/lib/pq/oid"
	"pqprobe/pkg/metadata/internal"
)

type (
	// Field contains table's fields properties.
	Field struct {
		ID              int64  `db:"attnum"`
		TypeID          uint32 `db:"atttypid"`
		Name            string `db:"attname"`
		Type            string `db:"format_type"`
		Nullable        bool   `db:"attnotnull"`
		IsPrimary       bool   `db:"indisprimary"`
		IsUnique        bool   `db:"indisunique"`
		IsValid         bool   `db:"indisvalid"`
		IndexDefinition string `db:"pg_get_indexdef"`
	}
)

// GoName returns field info in camelcase.
// Useful for struct field name.
func (t Field) GoName() string {
	return internal.Camelify(internal.VariableNameRule(t.Name))
}

// GoTypeString returns this field's type as equivalent Golang type.
// Useful for struct field type.
// See https://godoc.org/github.com/lib/pq#hdr-Data_Types for conventions.
func (t Field) GoTypeString() string {
	return internal.TypeString(oid.Oid(t.TypeID))
}
