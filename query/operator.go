package query

type Operator uint8

const (
	OperatorAnd            Operator = iota
	OperatorOr             Operator = iota
	OperatorEqual          Operator = iota
	OperatorNotEqual       Operator = iota
	OperatorGreater        Operator = iota
	OperatorLess           Operator = iota
	OperatorGreaterOrEqual Operator = iota
	OperatorLessOrEqual    Operator = iota
	OperatorIn             Operator = iota
)

type Condition struct {
	Operation Operator
	Left      interface{}
	Right     interface{}
}

func And(left *Condition, right *Condition) *Condition {
	return &Condition{
		Operation: OperatorAnd,
		Left:      left,
		Right:     right,
	}
}

func Or(left *Condition, right *Condition) *Condition {
	return &Condition{
		Operation: OperatorOr,
		Left:      left,
		Right:     right,
	}
}

func Equal(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorEqual,
		Left:      column,
		Right:     value,
	}
}

func NotEqual(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorNotEqual,
		Left:      column,
		Right:     value,
	}
}

func GreaterThan(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorGreater,
		Left:      column,
		Right:     value,
	}
}

func LessThan(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorLess,
		Left:      column,
		Right:     value,
	}
}

func GreaterOrEqual(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorGreaterOrEqual,
		Left:      column,
		Right:     value,
	}
}

func LessOrEqual(column string, value interface{}) *Condition {
	return &Condition{
		Operation: OperatorLessOrEqual,
		Left:      column,
		Right:     value,
	}
}

func In(column string, values []interface{}) *Condition {
	return &Condition{
		Operation: OperatorIn,
		Left:      column,
		Right:     values,
	}
}
