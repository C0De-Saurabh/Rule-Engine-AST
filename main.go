package main

import (
	"fmt"
	"strconv"
	"strings"
)

// NodeType defines the type of the AST node
type NodeType int

const (
	NodeTypeNumber NodeType = iota
	NodeTypeOperator
	NodeTypeBoolean
	NodeTypeIdentifier
	NodeTypeString
)

// ASTNode represents a node in the abstract syntax tree
type ASTNode struct {
	Type  NodeType
	Value string
	Left  *ASTNode
	Right *ASTNode
}

// Tokenize splits the input into tokens
func Tokenize(input string) []string {
	var tokens []string
	var token strings.Builder
	inString := false

	for _, char := range input {
		switch {
		case char == '\'':
			inString = !inString
			token.WriteRune(char)
			if !inString {
				tokens = append(tokens, token.String())
				token.Reset()
			}
		case char == ' ' && !inString:
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
		case char == '(' || char == ')':
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			tokens = append(tokens, string(char))
		default:
			token.WriteRune(char)
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}

	return tokens
}

// ParseExpression recursively parses the expression to build an AST
func ParseExpression(tokens []string) (*ASTNode, []string) {
	if len(tokens) == 0 {
		return nil, tokens
	}

	// Parse the left side
	left, remainingTokens := ParseTerm(tokens)

	if len(remainingTokens) == 0 {
		return left, remainingTokens
	}

	// Parse logical operators (AND, OR)
	op := remainingTokens[0]
	if op == "AND" || op == "OR" {
		remainingTokens = remainingTokens[1:]

		// Parse the right side
		right, remainingTokens := ParseExpression(remainingTokens)

		// Create the operator node with left and right children
		node := &ASTNode{
			Type:  NodeTypeOperator,
			Value: op,
			Left:  left,
			Right: right,
		}

		return node, remainingTokens
	}

	return left, remainingTokens
}

// ParseTerm parses a term (identifier, number, boolean, string, comparison, or sub-expression) in the expression
func ParseTerm(tokens []string) (*ASTNode, []string) {
	if len(tokens) == 0 {
		return nil, tokens
	}

	// Handle parentheses for sub-expressions
	if tokens[0] == "(" {
		expr, remainingTokens := ParseExpression(tokens[1:])
		if len(remainingTokens) > 0 && remainingTokens[0] == ")" {
			return expr, remainingTokens[1:]
		}
		return nil, tokens
	}

	// Handle comparison operators (> , < , =)
	if len(tokens) > 2 {
		if tokens[1] == ">" || tokens[1] == "<" || tokens[1] == "=" {
			left := &ASTNode{
				Type:  NodeTypeIdentifier,
				Value: tokens[0],
			}
			op := tokens[1]
			right := &ASTNode{
				Type:  NodeTypeIdentifier,
				Value: tokens[2],
			}

			remainingTokens := tokens[3:]
			if _, err := strconv.Atoi(tokens[2]); err == nil {
				right.Type = NodeTypeNumber
			} else if strings.HasPrefix(tokens[2], "'") && strings.HasSuffix(tokens[2], "'") {
				right.Type = NodeTypeString
			}

			return &ASTNode{
				Type:  NodeTypeOperator,
				Value: op,
				Left:  left,
				Right: right,
			}, remainingTokens
		}
	}

	token := tokens[0]
	tokens = tokens[1:]

	// Check if the token is a number
	if _, err := strconv.Atoi(token); err == nil {
		return &ASTNode{
			Type:  NodeTypeNumber,
			Value: token,
		}, tokens
	}

	// Check if the token is a boolean value (true/false)
	if token == "true" || token == "false" {
		return &ASTNode{
			Type:  NodeTypeBoolean,
			Value: token,
		}, tokens
	}

	// Check if the token is a string (e.g., 'Sales')
	if strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'") {
		return &ASTNode{
			Type:  NodeTypeString,
			Value: token,
		}, tokens
	}

	// Assume it's an identifier (e.g., age, department)
	return &ASTNode{
		Type:  NodeTypeIdentifier,
		Value: token,
	}, tokens
}

// PrintAST prints the AST in a readable format
func PrintAST(node *ASTNode, indent string) {
	if node == nil {
		return
	}

	fmt.Println(indent + node.Value)

	if node.Left != nil {
		PrintAST(node.Left, indent+"  ")
	}
	if node.Right != nil {
		PrintAST(node.Right, indent+"  ")
	}
}

func ConvertASTToExpression(node *ASTNode) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case NodeTypeNumber, NodeTypeBoolean, NodeTypeString, NodeTypeIdentifier:
		return node.Value
	case NodeTypeOperator:
		// Handle operators with parentheses
		left := ConvertASTToExpression(node.Left)
		right := ConvertASTToExpression(node.Right)

		if node.Value == "AND" || node.Value == "OR" {
			// Logical operators need to ensure proper parentheses
			return fmt.Sprintf("(%s %s %s)", left, node.Value, right)
		}
		// Comparison operators also need to ensure proper parentheses
		return fmt.Sprintf("%s %s %s", left, node.Value, right)
	default:
		return ""
	}
}

func main() {
	input := "((age > 30 AND department = 'Marketing')) AND (salary > 20000 OR experience > 5)"
	tokens := Tokenize(input)
	ast, _ := ParseExpression(tokens)

	fmt.Println("Abstract Syntax Tree:")
	PrintAST(ast, "")

	fmt.Println("Converted Back to Expression:")

	fmt.Println(ConvertASTToExpression(ast))
}
