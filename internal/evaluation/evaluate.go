package evaluation

import (
	"rule-engine/internal/ast"
	"strings"
)

// EvaluateNode evaluates the AST node against the provided data.
func EvaluateNode(node *ast.Node, data map[string]interface{}) bool {
	if node == nil {
		return false
	}

	// If the node has logical operators (AND, OR), evaluate the left and right nodes.
	switch node.Operator {
	case "and", "&&":
		return EvaluateNode(node.Left, data) && EvaluateNode(node.Right, data)
	case "or", "||":
		return EvaluateNode(node.Left, data) || EvaluateNode(node.Right, data)
	}

	// Leaf node: evaluate the condition.
	fieldValue, exists := data[node.Field]
	if !exists {
		return false
	}

	return compareValues(fieldValue, node.Value, node.Operator)
}

// compareValues compares the actual data value with the value in the AST node based on the operator.
func compareValues(actual, expected interface{}, operator string) bool {
	switch actual := actual.(type) {
	case int:
		expected, ok := expected.(int)
		if !ok {
			return false
		}
		switch operator {
		case ">":
			return actual > expected
		case "=":
			return actual == expected
		case "<":
			return actual < expected
		}

	case string:
		expected, ok := expected.(string)
		if !ok {
			return false
		}
		if operator == "=" {
			return strings.EqualFold(actual, expected) // Case-insensitive comparison
		}
	}
	return false
}
