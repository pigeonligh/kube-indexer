package dataprocessor

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

type visitor struct{}

func (visitor) Visit(node *ast.Node) {
	if node != nil {
		switch n := (*node).(type) {
		case *ast.MemberNode:
			newNode := &ast.CallNode{
				Callee:    &ast.IdentifierNode{Value: "_dot"},
				Arguments: []ast.Node{n.Node, n.Property},
			}
			ast.Patch(node, newNode)

		case *ast.SliceNode:
			from := n.From
			if from == nil {
				from = &ast.NilNode{}
			}
			to := n.To
			if to == nil {
				to = &ast.NilNode{}
			}

			newNode := &ast.CallNode{
				Callee:    &ast.IdentifierNode{Value: "_range"},
				Arguments: []ast.Node{n.Node, from, to},
			}
			ast.Patch(node, newNode)

		case *ast.BinaryNode:
			newNode := &ast.BinaryNode{
				Operator: n.Operator,
				Left: &ast.CallNode{
					Callee:    &ast.IdentifierNode{Value: "_val"},
					Arguments: []ast.Node{n.Left},
				},
				Right: &ast.CallNode{
					Callee:    &ast.IdentifierNode{Value: "_val"},
					Arguments: []ast.Node{n.Right},
				},
			}
			ast.Patch(node, newNode)
		}
	}
}

func EvalExpr(src Source, e string, valueMap map[string]Object) (Object, error) {
	env := make(map[string]any)
	for k, v := range valueMap {
		env[k] = v
	}
	env["_dot"] = evalDot(src)
	env["_range"] = evalRange(src)
	env["_val"] = evalValue(src)
	env["len"] = evalLen()

	options := make([]expr.Option, 0)
	options = append(options, expr.Env(env))
	options = append(options, expr.Patch(visitor{}))

	prog, err := expr.Compile(e, options...)
	if err != nil {
		return nil, err
	}

	result, err := expr.Run(prog, env)
	if err != nil {
		return nil, err
	}
	return NewObject(result), nil
}

func evalDot(src Source) any {
	return func(obj Object, key any) Object {
		obj = UnrefObject(src, obj)

		switch key := key.(type) {
		case int:
			obj = obj.GetIndex(key)

		case string:
			obj = obj.Get(key)

		default:
			return NewObject(nil)
		}

		obj = UnrefObject(src, obj)
		return NewObject(obj)
	}
}

func evalRange(src Source) any {
	return func(obj Object, f, t any) Object {
		obj = UnrefObject(src, obj)

		var fp, tp *int
		if f, ok := f.(int); ok {
			fp = &f
		}
		if t, ok := t.(int); ok {
			tp = &t
		}
		return obj.Range(fp, tp)
	}
}

func evalValue(src Source) any {
	return func(obj any) any {
		if obj, ok := obj.(Object); ok {
			return UnrefObject(src, obj).Value()
		}
		return obj
	}
}

func evalLen() any {
	return func(o Object) Object {
		return NewObject(o.Len())
	}
}
