package main

import (
	"fmt"
	"log"
	"net/http"
	"rule-engine/internal/api"
)

func main() {
	router := api.SetupRouter()

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
