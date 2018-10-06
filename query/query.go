package query

type Query struct {
	from       string
	columns    []string
	conditions []*Condition
	take       int
	skip       int
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) From(from string) *Query {
	q.from = from
	return q
}

func (q *Query) Columns(columns ...string) *Query {
	q.columns = columns
	return q
}

func (q *Query) Where(conditions ...*Condition) *Query {
	q.conditions = append(q.conditions, conditions...)
	return q
}

func (q *Query) Take(take int) *Query {
	q.take = take
	return q
}

func (q *Query) Skip(skip int) *Query {
	q.skip = skip
	return q
}

func (q Query) Build(dialectKind DialectKind) (string, []interface{}, error) {
	var d dialect
	switch dialectKind {
	case DialectKindCassandra:
		d = &cassandraDialect{}
	}
	if d == nil {
		return "", nil, ErrInvalidDialect
	}
	return d.BuildQuery(q)
}
