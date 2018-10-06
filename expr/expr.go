package expr

import (
	"go/ast"
	"go/parser"
)

// Expression
type Expression struct {
	Source string
	tree   ast.Expr
}

func New(source string) *Expression {
	return &Expression{
		Source: source,
	}
}

func (e *Expression) Parse() error {
	tree, err := parser.ParseExpr(e.Source)
	if err != nil {
		return err
	}
	e.tree = tree
	return nil
}

func (e *Expression) Evaluate() (interface{}, error) {
	if e.tree == nil {
		err := e.Parse()
		if err != nil {
			return false, err
		}
	}
	eval := newEvaluator()
	return eval.Evaluate(e.tree)
}

func (e *Expression) MustEvaluate() (interface{}) {
	res, err := e.Evaluate()
	if err != nil {
		panic(err)
	}
	return res
}
