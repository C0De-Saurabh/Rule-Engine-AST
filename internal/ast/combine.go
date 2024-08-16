package ast

// Range represents a range with a lower and upper bound.
type Range struct {
	Lower int
	Upper int
}

// CombineASTs combines an array of AST nodes while reducing redundancy.
func CombineASTs(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}

	result := nodes[0]
	for _, node := range nodes[1:] {
		result = combineNodes(result, node)
	}

	return result
}

// normalize converts AST nodes into a canonical form.
func normalize(node *Node) *Node {
	if node == nil {
		return nil
	}

	// Normalize child nodes recursively
	if node.Operator == "and" || node.Operator == "or" {
		node.Left = normalize(node.Left)
		node.Right = normalize(node.Right)
		return node
	}

	// For 'in' operator, ensure it's in normalized form
	if node.Operator == "in" {
		return &Node{
			Field:    node.Field,
			Operator: "in",
			Value:    node.Value,
		}
	}

	return node
}

// combineNodes combines two nodes with minimal redundancy.
func combineNodes(node1, node2 *Node) *Node {
	if node1 == nil {
		return node2
	}
	if node2 == nil {
		return node1
	}

	// Normalize both nodes before combining
	normalizedNode1 := normalize(node1)
	normalizedNode2 := normalize(node2)

	if normalizedNode1.Field == normalizedNode2.Field {
		return mergeConditions(normalizedNode1, normalizedNode2)
	}

	return &Node{
		Operator: "and",
		Left:     normalizedNode1,
		Right:    normalizedNode2,
	}
}

// mergeConditions merges two conditions on the same field.
func mergeConditions(node1, node2 *Node) *Node {
	return &Node{
		Field:    node1.Field,
		Operator: node1.Operator,
		Value:    combineValues(node1.Value, node2.Value),
	}
}

// combineValues combines values based on the operator.
func combineValues(value1, value2 interface{}) interface{} {
	switch v1 := value1.(type) {
	case int:
		if v2, ok := value2.(int); ok {
			return Range{
				Lower: min(v1, v2),
				Upper: max(v1, v2),
			}
		}
	case string:
		if v2, ok := value2.(string); ok {
			return []string{v1, v2}
		}
	}
	return value1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
