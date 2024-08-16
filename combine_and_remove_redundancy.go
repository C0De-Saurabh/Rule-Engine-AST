package main

import "errors"

// combineRules combines multiple ASTs into a single AST with an AND operator.
func combineRules(rules []string) (*ASTNode, error) {
	if len(rules) == 0 {
		return nil, errors.New("no rules to combine")
	}

	if len(rules) == 1 {
		return createRule(rules[0])
	}

	combinedAST := &ASTNode{
		Type:  "operator",
		Value: "AND",
	}

	var leftAST *ASTNode
	for i, rule := range rules {
		ast, err := createRule(rule)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			leftAST = ast
		} else {
			combinedAST.Left = leftAST
			combinedAST.Right = ast
			leftAST = combinedAST
		}
	}

	return combinedAST, nil
}
