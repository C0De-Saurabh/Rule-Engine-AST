package main

import (
	"errors"
	"fmt"
)

// evaluateRule evaluates the AST against provided JSON data
func evaluateRule(ast *Node, data map[string]interface{}) (bool, error) {
	if ast == nil {
		return false, errors.New("AST cannot be nil")
	}

	switch ast.Type {
	case "operator":
		// Evaluate the left and right subtrees
		leftEval, err := evaluateRule(ast.Left, data)
		if err != nil {
			return false, err
		}
		rightEval, err := evaluateRule(ast.Right, data)
		if err != nil {
			return false, err
		}

		// Apply the operator
		switch ast.Value {
		case "AND":
			return leftEval && rightEval, nil
		case "OR":
			return leftEval || rightEval, nil
		default:
			return false, fmt.Errorf("unsupported operator: %s", ast.Value)
		}

	case "operand":
		// Check if the operand matches the data
		value, ok := data[ast.Value]
		if !ok {
			return false, fmt.Errorf("missing attribute: %s", ast.Value)
		}

		// For simplicity, assume the operand is a boolean in the data
		return value.(bool), nil

	default:
		return false, fmt.Errorf("unknown node type: %s", ast.Type)
	}
}

// validateData checks that the data contains all necessary attributes
func validateData(data map[string]interface{}) error {
	if len(data) == 0 {
		return errors.New("data cannot be empty")
	}

	// Additional validation rules can be added here
	return nil
}
