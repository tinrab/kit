package query

type Update struct {
	table       string
	assignments []Assignment
	conditions  []*Condition
}

type Assignment struct {
	Column string
	Value  interface{}
}

func NewUpdate() *Update {
	return &Update{
		assignments: make([]Assignment, 0),
		conditions:  make([]*Condition, 0),
	}
}

func (u *Update) Table(table string) *Update {
	u.table = table
	return u
}

func (u *Update) Assign(column string, value interface{}) *Update {
	u.assignments = append(u.assignments, Assignment{
		Column: column,
		Value:  value,
	})
	return u
}

func (u *Update) Where(conditions ...*Condition) *Update {
	u.conditions = append(u.conditions, conditions...)
	return u
}

func (u Update) Build(dialectKind DialectKind) (string, []interface{}, error) {
	var d dialect
	switch dialectKind {
	case DialectKindCassandra:
		d = &cassandraDialect{}
	}
	if d == nil {
		return "", nil, ErrInvalidDialect
	}
	return d.BuildUpdate(u)
}
