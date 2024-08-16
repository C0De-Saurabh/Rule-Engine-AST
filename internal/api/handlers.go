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
		RuleID     string `json:"rule_id"`
		RuleString string `json:"rule_string"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Save the rule
	if err := storage.SaveRule(request.RuleID, request.RuleString); err != nil {
		http.Error(w, "Failed to save rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "Rule Saved successfully",
		"rule_id": request.RuleID,
	})
}

// CombineRules handles the combination of multiple rules.
func CombineRules(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RuleID string   `json:"rule_id"`
		Rules  []string `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(request.Rules) < 2 {
		http.Error(w, "At least two rules are required to combine", http.StatusBadRequest)
		return
	}

	// Combine the rules with the logical AND operator
	combinedRule := request.Rules[0]
	for i := 1; i < len(request.Rules); i++ {
		combinedRule = combinedRule + " && " + request.Rules[i]
	}

	// Save the combined rule
	if err := storage.SaveRule(request.RuleID, combinedRule); err != nil {
		http.Error(w, "Failed to save combined rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "Rules combined successfully",
		"rule_id": request.RuleID,
	})
}

// EvaluateRequest defines the structure of the request body for evaluating rules.
type EvaluateRequest struct {
	RuleID string                 `json:"rule_id"`
	Data   map[string]interface{} `json:"data"`
}

// EvaluateRule handles the evaluation of a rule against provided data.
func EvaluateRule(w http.ResponseWriter, r *http.Request) {
	var request EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ruleString, err := storage.RetrieveRule(request.RuleID)
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
			data[key] = int(v)
		default:
			data[key] = v
		}
	}

	log.Printf("Deserialized data: %+v", data)

	astNode := ast.ParseAST(ruleString)
	if astNode == nil {
		http.Error(w, "Failed to parse rule string", http.StatusBadRequest)
		return
	}

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = evaluation.EvaluateNode(astNode, data)
		log.Printf("Evaluation Result: %v", result)
	}()

	wg.Wait()

	response := struct {
		Result bool   `json:"result"`
		Rule   string `json:"rule"`
	}{
		Result: result,
		Rule:   ruleString,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
