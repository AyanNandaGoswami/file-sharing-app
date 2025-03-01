package main

import (
	"fmt"
	"net/http"

	"github.com/AyanNandaGoswami/microservices/file-sharing-app/authentication/internal/api/routes"
)

func main() {
	// Initialize routes
	routes.InitializeRoutes()

	// Start HTTP server on port 4001
	fmt.Println("Authenticaton service is listening on port 4001...")

	// ListenAndServe starts an HTTP server with a given address and handler.
	// If the address is blank, ":http" is used (i.e., "localhost:8080").
	if err := http.ListenAndServe(":4001", nil); err != nil {
		// Error starting server, print error and exit
		fmt.Println("Error starting server:", err)
	}
}
