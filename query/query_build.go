package query

import (
	"fmt"
	"strings"
)

func (q Query) Build(d Dialect) (string, []interface{}, error) {
	switch d {
	case DialectCassandra:
		return q.buildCassandra()
	}
	return "", nil, ErrInvalidDialect
}

func (q Query) buildCassandra() (string, []interface{}, error) {
	b := strings.Builder{}
	args := make([]interface{}, 0)

	b.WriteString("SELECT ")
	for i, c := range q.columns {
		b.WriteString(c)
		if i < len(q.columns)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(" FROM ")
	b.WriteString(q.from)

	b.WriteString(" WHERE ")
	for i, c := range q.conditions {
		cs, ca, err := q.buildCondition(c)
		if err != nil {
			return "", nil, err
		}
		b.WriteString(cs)
		args = append(args, ca...)

		if i < len(q.conditions)-1 {
			b.WriteString(" AND ")
		}
	}

	if q.take != -1 {
		b.WriteString(" LIMIT ")
		b.WriteString(fmt.Sprint(q.take))
	}

	b.WriteString(" ALLOW FILTERING")

	return b.String(), args, nil
}

func (q Query) buildCondition(c *Condition) (string, []interface{}, error) {
	stmt := strings.Builder{}
	args := make([]interface{}, 0)

	if lc, ok := c.Left.(*Condition); ok {
		ls, la, err := q.buildCondition(lc)
		if err != nil {
			return "", nil, err
		}
		stmt.WriteString(ls)
		args = append(args, la...)
	} else if lc, ok := c.Left.(string); ok {
		stmt.WriteString(lc)
	} else {
		return "", nil, ErrInvalidCondition
	}

	switch c.Operation {
	case OperatorAnd:
		stmt.WriteString(" AND ")
	case OperatorOr:
		stmt.WriteString(" OR ")
	case OperatorEqual:
		stmt.WriteString(" = ")
	case OperatorNotEqual:
		stmt.WriteString(" != ")
	case OperatorGreater:
		stmt.WriteString(" > ")
	case OperatorLess:
		stmt.WriteString(" < ")
	case OperatorGreaterOrEqual:
		stmt.WriteString(" >= ")
	case OperatorLessOrEqual:
		stmt.WriteString(" <= ")
	}

	if rc, ok := c.Right.(*Condition); ok {
		rs, ra, err := q.buildCondition(rc)
		if err != nil {
			return "", nil, err
		}
		stmt.WriteString(rs)
		args = append(args, ra...)
	} else {
		stmt.WriteString("?")
		args = append(args, c.Right)
	}

	return stmt.String(), args, nil
}
