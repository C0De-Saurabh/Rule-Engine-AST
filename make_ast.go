package main

import (
	"errors"
	"strings"
)

// ASTNode represents a node in the Abstract Syntax Tree (AST).
type ASTNode struct {
	Type  string   `json:"type"`
	Value string   `json:"value"`
	Left  *ASTNode `json:"left,omitempty"`
	Right *ASTNode `json:"right,omitempty"`
}

// createRule converts a rule string into an AST.
func createRule(rule string) (*ASTNode, error) {
	tokens := strings.Fields(rule)
	if len(tokens) < 3 {
		return nil, errors.New("invalid rule")
	}

	root := &ASTNode{
		Type:  "operator",
		Value: tokens[1],
	}

	root.Left = &ASTNode{
		Type:  "operand",
		Value: tokens[0],
	}

	root.Right = &ASTNode{
		Type:  "operand",
		Value: tokens[2],
	}

	return root, nil
}
