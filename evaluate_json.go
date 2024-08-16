package main

import (
	"errors"
	"fmt"
)

// evaluateRule evaluates the AST against provided JSON data.
func evaluateRule(ast *ASTNode, data map[string]interface{}) (bool, error) {
	if ast == nil {
		return false, errors.New("invalid AST")
	}

	switch ast.Type {
	case "operator":
		leftValue, err := evaluateRule(ast.Left, data)
		if err != nil {
			return false, err
		}
		rightValue, err := evaluateRule(ast.Right, data)
		if err != nil {
			return false, err
		}

		switch ast.Value {
		case "AND":
			return leftValue && rightValue, nil
		case "OR":
			return leftValue || rightValue, nil
		case ">":
			return leftValue && rightValue, nil
		case "==":
			return leftValue == rightValue, nil
		default:
			return false, fmt.Errorf("unsupported operator: %s", ast.Value)
		}

	case "operand":
		value, ok := data[ast.Value]
		if !ok {
			return false, fmt.Errorf("key %s not found in data", ast.Value)
		}
		return value == ast.Value, nil

	default:
		return false, fmt.Errorf("unknown AST node type: %s", ast.Type)
	}
}
