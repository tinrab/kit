package query

import (
	"strings"
)

func (u Update) Build(d Dialect) (string, []interface{}, error) {
	switch d {
	case DialectCassandra:
		return u.buildCassandra()
	}
	return "", nil, ErrInvalidDialect
}

func (u Update) buildCassandra() (string, []interface{}, error) {
	b := strings.Builder{}
	args := make([]interface{}, 0)

	b.WriteString("UPDATE ")
	b.WriteString(u.table)

	b.WriteString(" SET ")
	for j, a := range u.assignments {
		b.WriteString(a.Column)
		b.WriteString(" = ?")
		if j < len(u.assignments)-1 {
			b.WriteString(", ")
		}
		args = append(args, a.Value)
	}

	b.WriteString(" WHERE ")
	for i, c := range u.conditions {
		cs, ca, err := u.buildCondition(c)
		if err != nil {
			return "", nil, err
		}
		b.WriteString(cs)
		args = append(args, ca...)

		if i < len(u.conditions)-1 {
			b.WriteString(" AND ")
		}
	}

	return b.String(), args, nil
}

func (u Update) buildCondition(c *Condition) (string, []interface{}, error) {
	stmt := strings.Builder{}
	args := make([]interface{}, 0)

	if lc, ok := c.Left.(*Condition); ok {
		ls, la, err := u.buildCondition(lc)
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
		rs, ra, err := u.buildCondition(rc)
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
