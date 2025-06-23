package main

import (
	"time"

	"github.com/aditya-goyal-omniful/ims/docs"
	localConfig "github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/aditya-goyal-omniful/ims/pkg/middlewares"
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

	// Swagger metadata
	docs.SwaggerInfo.Title = "Inventory Management Service"
	docs.SwaggerInfo.Description = "API documentation for IMS"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8087"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	
	// Initialize HTTP server
	server := http.InitializeServer(
		":"+port,
		10*time.Second,
		10*time.Second,
		70*time.Second,
		false,
	)
	
	server.Use(config.Middleware())
	server.Use(middlewares.RequestLogger())

	routes.SetupRoutes(server)

	log.Infof("Starting server on port", port)
	if err := server.StartServer("ims"); err != nil {
		log.Panic("Failed to start server: %v", err)
	}
}