package ast

import (
	"log"
	"strings"
)

// PrintAST prints the AST in a readable format.
func PrintAST(node *Node, indent string) {
	if node == nil {
		return
	}

	// Print the current node
	if node.Field != "" {
		// Leaf node (comparison operation)
		log.Printf("%sField: %s, Operator: %s, Value: %v\n", indent, node.Field, node.Operator, node.Value)
	} else {
		// Non-leaf node (AND/OR operation)
		log.Printf("%sOperator: %s\n", indent, strings.ToUpper(node.Operator))
	}

	// Recur on the left and right children with increased indentation
	if node.Left != nil {
		PrintAST(node.Left, indent+"  ")
	}
	if node.Right != nil {
		PrintAST(node.Right, indent+"  ")
	}
}
