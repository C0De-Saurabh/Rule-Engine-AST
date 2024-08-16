package ast

// Node struct for my Syntax Tree
type Node struct {
	Field    string
	Operator string
	Value    interface{}
	Left     *Node
	Right    *Node
}

//NewNode creates a new Node

func NewNode(field string, operator string, value interface{}) *Node {
	return &Node{Field: field, Operator: operator, Value: value}
}
