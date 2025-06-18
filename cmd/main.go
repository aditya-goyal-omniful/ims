package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aditya-goyal-omniful/ims/pkg/config"
	"github.com/aditya-goyal-omniful/ims/pkg/routes"
	"github.com/omniful/go_commons/http"
)

func main() {
	config.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	server := http.InitializeServer(
		":"+port,
		10*time.Second,  // Read timeout
		10*time.Second,  // Write timeout
		70*time.Second,  // Idle timeout
		false,
	)

	routes.SetupRoutes(server)

	fmt.Println("Starting server on port", port)
	if err := server.StartServer("ims"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
