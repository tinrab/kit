package query

import (
	"strings"
)

func (i Insert) Build(d Dialect) (string, []interface{}, error) {
	switch d {
	case DialectCassandra:
		return i.buildCassandra()
	}
	return "", nil, ErrInvalidDialect
}

func (i Insert) buildCassandra() (string, []interface{}, error) {
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
