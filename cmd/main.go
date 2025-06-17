package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	r := gin.Default()

	// Initialize routes
	// routes.SetupRoutes(r)

	r.Run()
}