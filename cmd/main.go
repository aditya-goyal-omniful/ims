package main

import (
	"time"

	localConfig "github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/aditya-goyal-omniful/ims/pkg/routes"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

func main() {
	// Initialize configuration from CONFIG_SOURCE
	if err := config.Init(15 * time.Second); err != nil {
		log.Panic("Failed to initialize config: %v", err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		log.Panic("Failed to get config context: %v", err)
	}

	port := config.GetString(ctx, "server.port")
	if port == "" {
		port = "8087"
	}

	localConfig.InitDB(ctx)
	localConfig.InitRedis()
	defer localConfig.RedisClient.Close()

	// Initialize HTTP server
	server := http.InitializeServer(
		":"+port,
		10*time.Second,
		10*time.Second,
		70*time.Second,
		false,
	)

	server.Use(config.Middleware())

	routes.SetupRoutes(server)

	log.Infof("Starting server on port", port)
	if err := server.StartServer("ims"); err != nil {
		log.Panic("Failed to start server: %v", err)
	}
}