package main

import (
	"errors"
	"strings"
)

// Node represents a single node in the AST
type Node struct {
	Type  string // "operator" or "operand"
	Left  *Node  // Left child node
	Right *Node  // Right child node
	Value string // Value for operand nodes
}

// createRule parses a rule string and returns the root node of the AST
func createRule(ruleString string) (*Node, error) {
	// Validate the rule string
	if err := validateRule(ruleString); err != nil {
		return nil, err
	}

	// Trim spaces and split the rule into tokens
	tokens := strings.Fields(ruleString)

	// Build the AST from tokens
	root, err := parseTokens(tokens)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// validateRule checks the validity of the rule string
func validateRule(ruleString string) error {
	if len(ruleString) == 0 {
		return errors.New("rule string cannot be empty")
	}

	// Add more validation rules as needed (e.g., balanced parentheses)
	return nil
}

// parseTokens recursively builds the AST from a list of tokens
func parseTokens(tokens []string) (*Node, error) {
	if len(tokens) == 0 {
		return nil, errors.New("no tokens to parse")
	}

	// Example for a simple implementation: parse binary operations
	if len(tokens) == 3 {
		// Simple binary operation: operand1 operator operand2
		node := &Node{
			Type:  "operator",
			Value: tokens[1],
			Left:  &Node{Type: "operand", Value: tokens[0]},
			Right: &Node{Type: "operand", Value: tokens[2]},
		}
		return node, nil
	}

	// Handle more complex cases (e.g., nested operations)
	// This will likely require a stack-based approach or recursive descent parsing

	return nil, errors.New("unable to parse tokens")
}
