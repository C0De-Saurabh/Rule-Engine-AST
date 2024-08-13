package main

import (
	"encoding/json"
	"net/http"
)

// handleCreateRule handles the creation of a new rule via the API
func handleCreateRule(w http.ResponseWriter, r *http.Request) {
	var ruleString string
	if err := json.NewDecoder(r.Body).Decode(&ruleString); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the rule string
	if err := validateRule(ruleString); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the AST
	ast, err := createRule(ruleString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the AST as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ast)
}

// handleCombineRules handles the combination of multiple rules via the API
func handleCombineRules(w http.ResponseWriter, r *http.Request) {
	var ruleStrings []string
	if err := json.NewDecoder(r.Body).Decode(&ruleStrings); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create ASTs for each rule string
	var asts []*Node
	for _, ruleString := range ruleStrings {
		ast, err := createRule(ruleString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		asts = append(asts, ast)
	}

	// Combine the ASTs
	combinedAST, err := combineRules(asts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the combined AST as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedAST)
}

// handleEvaluateRule handles the evaluation of a rule against provided data via the API
func handleEvaluateRule(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		AST  *Node                  `json:"ast"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the data
	if err := validateData(payload.Data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Evaluate the rule
	result, err := evaluateRule(payload.AST, payload.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the evaluation result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"result": result})
}
