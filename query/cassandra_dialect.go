package query

import (
	"fmt"
	"strings"
)

type cassandraDialect struct {
}

func (d *cassandraDialect) BuildQuery(q Query) (string, []interface{}, error) {
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
		cs, ca, err := d.buildCondition(c)
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

func (d *cassandraDialect) BuildInsert(i Insert) (string, []interface{}, error) {
	b := strings.Builder{}
	args := make([]interface{}, 0)

	b.WriteString("INSERT INTO ")
	b.WriteString(i.into)
	b.WriteString("(")
	for j, c := range i.columns {
		b.WriteString(c)
		if j < len(i.columns)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(")")

	b.WriteString(" VALUES ")

	for j, row := range i.rows {
		b.WriteString("(")
		for k, v := range row {
			b.WriteString("?")
			args = append(args, v)
			if k < len(row)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
		if j < len(i.rows)-1 {
			b.WriteString(", ")
		}
	}

	return b.String(), args, nil
}

func (d *cassandraDialect) BuildUpdate(u Update) (string, []interface{}, error) {
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
		cs, ca, err := d.buildCondition(c)
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

func (d *cassandraDialect) buildCondition(c *Condition) (string, []interface{}, error) {
	stmt := strings.Builder{}
	args := make([]interface{}, 0)

	if lc, ok := c.Left.(*Condition); ok {
		ls, la, err := d.buildCondition(lc)
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
	case OperatorIn:
		stmt.WriteString(" IN ")
	}

	if rc, ok := c.Right.(*Condition); ok {
		rs, ra, err := d.buildCondition(rc)
		if err != nil {
			return "", nil, err
		}
		stmt.WriteString(rs)
		args = append(args, ra...)
	} else if values, ok := c.Right.([]interface{}); ok {
		stmt.WriteString("(")
		for i := 0; i < len(values); i++ {
			stmt.WriteString("?")
			if i < len(values)-1 {
				stmt.WriteString(", ")
			}
		}
		stmt.WriteString(")")
		args = append(args, values...)
	} else {
		stmt.WriteString("?")
		args = append(args, c.Right)
	}

	return stmt.String(), args, nil
}
