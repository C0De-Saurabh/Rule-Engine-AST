package main

import (
	"encoding/json"
	"net/http"
)

// Struct for the request payload when creating a rule
type CreateRuleRequest struct {
	Rule string `json:"rule"`
}

// Struct for the request payload when combining rules
type CombineRulesRequest struct {
	Rules []string `json:"rules"`
}

// Struct for the request payload when evaluating a rule against JSON data
type EvaluateRuleRequest struct {
	AST  ASTNode                `json:"ast"`
	Data map[string]interface{} `json:"data"`
}

// Handler to create a rule and return its AST
func createRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the function to create the AST from the rule string
	ast, err := createRule(req.Rule)
	if err != nil {
		http.Error(w, "Failed to create rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the AST as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ast)
}

// Handler to combine multiple rules into a single AST
func combineRulesHandler(w http.ResponseWriter, r *http.Request) {
	var req CombineRulesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the function to combine the rules into one AST
	combinedAST, err := combineRules(req.Rules)
	if err != nil {
		http.Error(w, "Failed to combine rules: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the combined AST as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedAST)
}

// Handler to evaluate a rule against JSON data
func evaluateRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req EvaluateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the function to evaluate the AST against the provided data
	result, err := evaluateRule(&req.AST, req.Data)
	if err != nil {
		http.Error(w, "Failed to evaluate rule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the evaluation result as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"result": result})
}
