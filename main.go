package main

import (
	"log"
	"net/http"
)

func main() {
	// Define the API routes and their respective handlers
	http.HandleFunc("/create_rule", createRuleHandler)
	http.HandleFunc("/combine_rules", combineRulesHandler)
	http.HandleFunc("/evaluate_rule", evaluateRuleHandler)

	// Start the server on port 8080
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
