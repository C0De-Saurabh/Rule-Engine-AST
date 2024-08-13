package main

import (
	"errors"
)

// combineRules takes multiple ASTs and combines them into a single AST
func combineRules(rules []*Node) (*Node, error) {
	if len(rules) == 0 {
		return nil, errors.New("no rules provided")
	}

	// Combine ASTs by connecting them with an AND operator
	combinedRoot := rules[0]
	for _, rule := range rules[1:] {
		combinedRoot = &Node{
			Type:  "operator",
			Value: "AND",
			Left:  combinedRoot,
			Right: rule,
		}
	}

	// Optimize the combined AST
	combinedRoot = optimizeAST(combinedRoot)
	return combinedRoot, nil
}

// optimizeAST removes redundancy and simplifies the AST
func optimizeAST(root *Node) *Node {
	if root == nil {
		return nil
	}

	// Example: Remove redundant AND/OR nodes
	if root.Type == "operator" && root.Left != nil && root.Right != nil {
		if root.Value == root.Left.Value && root.Value == root.Right.Value {
			return root.Left
		}
	}

	// Recursively optimize the left and right subtrees
	root.Left = optimizeAST(root.Left)
	root.Right = optimizeAST(root.Right)
	return root
}
