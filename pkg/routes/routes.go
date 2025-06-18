package routes

import (
	"github.com/aditya-goyal-omniful/ims/pkg/controllers"
	commonHttp "github.com/omniful/go_commons/http"
)

func SetupRoutes(server *commonHttp.Server) {
	// Hub routes
	server.GET("/hubs", controllers.GetHubs)
	server.GET("/hubs/:id", controllers.GetHubByID)
	server.POST("/hubs", controllers.CreateHub)
	server.DELETE("/hubs/:id", controllers.DeleteHub)
	server.PUT("/hubs/:id", controllers.UpdateHub)

	// SKU routes
	server.GET("/skus", controllers.GetSkus)
	server.GET("/skus/:id", controllers.GetSkuByID)
	server.POST("/skus", controllers.CreateSku)
	server.DELETE("/skus/:id", controllers.DeleteSku)
	server.PUT("/skus/:id", controllers.UpdateSku)

	// Inventory routes
	server.GET("/inventories", controllers.GetInventories)
	server.GET("/inventories/:id", controllers.GetInventoryByID)
	server.POST("/inventories", controllers.CreateInventory)
	server.DELETE("/inventories/:id", controllers.DeleteInventory)
	server.PUT("/inventories/:id", controllers.UpdateInventory)
}
