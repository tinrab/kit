package query

type DialectKind uint8

const (
	DialectKindCassandra DialectKind = iota
)

type dialect interface {
	BuildQuery(q Query) (string, []interface{}, error)
	BuildInsert(i Insert) (string, []interface{}, error)
	BuildUpdate(u Update) (string, []interface{}, error)
}
