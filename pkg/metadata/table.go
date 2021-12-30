package metadata

import "pqprobe/pkg/metadata/internal"

type (
	// Relation describes relation within database.
	Relation struct {
		Schema string `db:"schema"`
		Name   string `db:"name"`
		Type   string `db:"type"`
		Owner  string `db:"owner"`
	}

	// Table contains Postgres table information.
	Table struct {
		Name           string
		Fields         []Field
		PrimaryKeys    []Field
		NonPrimaryKeys []Field
	}
)

// GoName returns table name in camelcase.
// Useful for struct field name.
func (t Table) GoName() string {
	return internal.Camelify(internal.VariableNameRule(t.Name))
}
