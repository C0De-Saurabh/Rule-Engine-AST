package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the application
	if err := initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Set up the routes
	setupRoutes()

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initialize() error {
	// Initialize configuration settings (if any)
	// Initialize database connection (if any)
	// Any other initial setup

	log.Println("Application initialized successfully")
	return nil
}

func setupRoutes() {
	// Route handlers for API endpoints
	http.HandleFunc("/create_rule", handleCreateRule)
	http.HandleFunc("/combine_rules", handleCombineRules)
	http.HandleFunc("/evaluate_rule", handleEvaluateRule)
}
