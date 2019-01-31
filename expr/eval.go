package expr

import (
	"github.com/pkg/errors"
	"go/ast"
	"go/token"
	"math"
	"strconv"
)

type EvaluatorFunc func(s *Stack) interface{}

type evaluator struct {
	identifiers map[string]interface{}
	stack       *Stack
	err         error
}

func newEvaluator() *evaluator {
	e := &evaluator{
		identifiers: make(map[string]interface{}),
		stack:       NewStack(),
	}
	e.useCoreLibrary()
	return e
}

func (e *evaluator) Evaluate(node ast.Node) (interface{}, error) {
	ast.Walk(e, node)
	if e.err != nil {
		return nil, e.err
	}
	if !e.stack.IsEmpty() {
		return e.stack.Tail(), nil
	}
	return nil, errors.New("empty expression")
}

func (e *evaluator) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BasicLit:
		e.VisitBasicLiteral(n)
	case *ast.Ident:
		e.VisitIdentifier(n)
	case *ast.BinaryExpr:
		e.VisitBinaryExpression(n)
	case *ast.UnaryExpr:
		e.VisitUnaryExpression(n)
	case *ast.ParenExpr:
		e.VisitParenthesesExpression(n)
	case *ast.CallExpr:
		e.VisitCallExpression(n)
	default:
	}
	return nil
}

func (e *evaluator) VisitBasicLiteral(n *ast.BasicLit) {
	var x interface{}
	switch n.Kind {
	case token.INT, token.FLOAT:
		x, _ = strconv.ParseFloat(n.Value, 64)
	case token.STRING:
		x = n.Value
	}
	e.stack.Push(x)
}

func (e *evaluator) VisitBinaryExpression(n *ast.BinaryExpr) {
	ast.Walk(e, n.X)
	ast.Walk(e, n.Y)

	y := e.stack.Pop()
	x := e.stack.Pop()

	switch x := x.(type) {
	case float64:
		switch y := y.(type) {
		case float64:
			switch n.Op {
			case token.ADD:
				e.stack.Push(x + y)
			case token.SUB:
				e.stack.Push(x - y)
			case token.MUL:
				e.stack.Push(x * y)
			case token.QUO:
				e.stack.Push(x / y)
			case token.EQL:
				e.stack.Push(x == y)
			case token.NEQ:
				e.stack.Push(x != y)
			default:
				e.err = errors.Errorf("invalid operation '%s' at [%v]", n.Op, n.Pos())
				return
			}
		}
	case string:
		switch y := y.(type) {
		case string:
			e.stack.Push(x + y)
		}
	}
}

func (e *evaluator) VisitUnaryExpression(n *ast.UnaryExpr) {
	ast.Walk(e, n.X)

	x := e.stack.Pop()

	switch x := x.(type) {
	case float64:
		switch n.Op {
		case token.SUB:
			e.stack.Push(-x)
		default:
			e.err = errors.Errorf("invalid operation '%s' at [%v]", n.Op, n.Pos())
			return
		}
	case bool:
		switch n.Op {
		case token.NOT:
			e.stack.Push(!x)
		default:
			e.err = errors.Errorf("invalid operation '%s' at [%v]", n.Op, n.Pos())
		}
	}
}

func (e *evaluator) VisitParenthesesExpression(n *ast.ParenExpr) {
	ast.Walk(e, n.X)
}

func (e *evaluator) VisitCallExpression(n *ast.CallExpr) {
	for _, arg := range n.Args {
		ast.Walk(e, arg)
	}
	ast.Walk(e, n.Fun)
}

func (e *evaluator) VisitIdentifier(n *ast.Ident) {
	if x, ok := e.identifiers[n.Name]; ok {
		switch x := x.(type) {
		case EvaluatorFunc:
			y := x(e.stack)
			e.stack.Push(y)
		default:
			e.stack.Push(x)
		}
	} else {
		e.err = errors.Errorf("unknown identifier '%s' at [%v]", n.Name, n.Pos())
	}
}

func (e *evaluator) RegisterFunc(name string, f EvaluatorFunc) {
	e.identifiers[name] = f
}

func (e *evaluator) useCoreLibrary() {
	e.identifiers["true"] = true
	e.identifiers["false"] = false
	e.identifiers["Pi"] = math.Pi
	e.identifiers["E"] = math.E

	e.RegisterFunc("floor", func(s *Stack) interface{} {
		return math.Floor(s.Pop().(float64))
	})
	e.RegisterFunc("ceil", func(s *Stack) interface{} {
		return math.Ceil(s.Pop().(float64))
	})
	e.RegisterFunc("round", func(s *Stack) interface{} {
		return math.Round(s.Pop().(float64))
	})
	e.RegisterFunc("pow", func(s *Stack) interface{} {
		y := s.Pop().(float64)
		x := s.Pop().(float64)
		return math.Pow(x, y)
	})
	e.RegisterFunc("sin", func(s *Stack) interface{} {
		return math.Sin(s.Pop().(float64))
	})
	e.RegisterFunc("cos", func(s *Stack) interface{} {
		return math.Cos(s.Pop().(float64))
	})
}
