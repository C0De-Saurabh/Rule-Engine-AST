package api

import (
	"github.com/gorilla/mux"
)

// SetupRouter initializes the router with routes and handlers.
func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/create_rule", CreateRule).Methods("POST")
	r.HandleFunc("/combine_rules", CombineRules).Methods("POST")
	r.HandleFunc("/evaluate_rule", EvaluateRule).Methods("POST")
	return r
}
