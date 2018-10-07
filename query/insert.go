package query

type Insert struct {
	into    string
	columns []string
	rows    [][]interface{}
}

func NewInsert() *Insert {
	return &Insert{
		columns: make([]string, 0),
		rows:    make([][]interface{}, 0),
	}
}

func (i *Insert) Into(into string) *Insert {
	i.into = into
	return i
}

func (i *Insert) Columns(columns ...string) *Insert {
	i.columns = columns
	return i
}

func (i *Insert) Rows(rows ...[]interface{}) *Insert {
	i.rows = rows
	return i
}

func (i *Insert) Row(values ...interface{}) *Insert {
	i.rows = append(i.rows, values)
	return i
}

func (i Insert) Build(dialectKind DialectKind) (string, []interface{}, error) {
	var d dialect
	switch dialectKind {
	case DialectKindCassandra:
		d = &cassandraDialect{}
	}
	if d == nil {
		return "", nil, ErrInvalidDialect
	}
	return d.BuildInsert(i)
}
