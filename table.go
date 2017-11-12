package pqprobe

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

	// Table contains Posgres table information.
	Table struct {
		Name        string
		Fields      []Field
		PrimaryKeys []Field
	}
)

// GoName returns table name in camelcase.
// Useful for struct field name.
func (t Table) GoName() string {
	return camelify(variableNameRule(t.Name))
}
