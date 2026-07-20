package main

import (
	"errors"
	"fmt"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
)

// ZeroDivVisitor проверяет AST-дерево на деление на ноль
type ZeroDivVisitor struct {
	err error
}

func (v *ZeroDivVisitor) Visit(node *ast.Node) {
	if v.err != nil {
		return
	}

	if binary, ok := (*node).(*ast.BinaryNode); ok && binary.Operator == "/" {
		if isConstantZero(binary.Right) {
			v.err = errors.New("Cannot divide by zero")
			return
		}
		if val, err := evalStaticNode(binary.Right); err == nil && val == 0 {
			v.err = errors.New("Cannot divide by zero")
			return
		}
	}
}

func isConstantZero(node ast.Node) bool {
	if intNode, ok := node.(*ast.IntegerNode); ok && intNode.Value == 0 {
		return true
	}
	if floatNode, ok := node.(*ast.FloatNode); ok && floatNode.Value == 0.0 {
		return true
	}
	return false
}

func evalStaticNode(node ast.Node) (float64, error) {
	exprStr := fmt.Sprintf("%v", node)
	program, err := expr.Compile(exprStr)
	if err != nil {
		return 0, err
	}
	res, err := expr.Run(program, nil)
	if err != nil {
		return 0, err
	}

	switch val := res.(type) {
	case int:
		return float64(val), nil
	case float64:
		return val, nil
	}
	return 0, errors.New("not a numeric constant")
}