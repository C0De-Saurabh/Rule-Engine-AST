package api

import (
	"encoding/json"
	"log"
	"net/http"
	"rule-engine/internal/ast"
	"rule-engine/internal/evaluation"
	"rule-engine/internal/storage"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateRule handles the creation of a rule from a string.
func CreateRule(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RuleID      string   `json:"rule_id"`
		RuleStrings []string `json:"rule_strings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Combine the rule strings into a single AST
	var astNodes []*ast.Node
	for _, ruleString := range request.RuleStrings {
		astNode := ast.ParseAST(ruleString)
		if astNode == nil {
			http.Error(w, "Failed to parse rule string", http.StatusBadRequest)
			return
		}
		astNodes = append(astNodes, astNode)
	}

	astRoot := ast.CombineASTs(astNodes)
	if astRoot == nil {
		http.Error(w, "Failed to combine ASTs", http.StatusInternalServerError)
		return
	}

	// Save the combined AST to MongoDB
	err := storage.SaveRule(request.RuleID, request.RuleStrings)
	if err != nil {
		http.Error(w, "Failed to save rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "rules combined and saved successfully",
		"rule_id": request.RuleID,
	})
}

// CombineRules handles the combination of multiple rules.
func CombineRules(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Rules []string `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(request.Rules) < 2 {
		http.Error(w, "At least two rules are required to combine", http.StatusBadRequest)
		return
	}

	// Combine the rules
	var astNodes []*ast.Node
	for _, ruleString := range request.Rules {
		astNode := ast.ParseAST(ruleString)
		if astNode == nil {
			http.Error(w, "Failed to parse rule string", http.StatusBadRequest)
			return
		}
		astNodes = append(astNodes, astNode)
	}

	astRoot := ast.CombineASTs(astNodes)
	if astRoot == nil {
		http.Error(w, "Failed to combine ASTs", http.StatusInternalServerError)
		return
	}
	log.Println("Combined AST Root:")
	// ast.PrintAST(astRoot, "",)

	// You might want to save or return the combined rule as needed
	// For now, let's just respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "rules combined successfully",
	})
}

// EvaluateRequest defines the structure of the request body for evaluating rules.
type EvaluateRequest struct {
	RuleID string                 `json:"rule_id"`
	Data   map[string]interface{} `bson:"data"`
}

func EvaluateRule(w http.ResponseWriter, r *http.Request) {
	var request EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	RuleStrings, err := storage.RetrieveRule(request.RuleID)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve rule", http.StatusInternalServerError)
		return
	}

	// Convert numeric values from float64 to int if needed
	data := make(map[string]interface{})
	for key, value := range request.Data {
		switch v := value.(type) {
		case float64:
			// Assuming that all numeric values should be int, cast them
			data[key] = int(v)
		default:
			// For non-numeric values, keep them as they are
			data[key] = v
		}
	}

	dataOnly := data

	log.Printf("Deserialized data: %+v", dataOnly)

	astNode := ast.ParseAST(RuleStrings[0])
	if astNode == nil {
		http.Error(w, "Failed to parse rule string", http.StatusBadRequest)
		return
	}

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = evaluation.EvaluateNode(astNode, dataOnly)
		log.Printf("Evaluation Result: %v", result)
	}()

	wg.Wait()

	response := struct {
		Result bool   `json:"result"`
		Rule   string `json:"rule"`
	}{
		Result: result,
		Rule:   RuleStrings[0],
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
