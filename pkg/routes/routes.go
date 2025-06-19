package routes

import (
	"github.com/aditya-goyal-omniful/ims/pkg/controllers"
	"github.com/aditya-goyal-omniful/ims/pkg/middlewares"
	commonHttp "github.com/omniful/go_commons/http"
)

func SetupRoutes(server *commonHttp.Server) {
	// Hub routes
	server.Group("/hubs", middlewares.AuthMiddleware(false)).
		GET("", controllers.GetHubs).
		GET("/:id", controllers.GetHubByID).
		POST("", controllers.CreateHub).
		DELETE("/:id", controllers.DeleteHub).
		PUT("/:id", controllers.UpdateHub)

	// SKU routes (Tenant + Seller)
	server.Group("/skus", middlewares.AuthMiddleware(true)).
		GET("", controllers.GetSkus).
		GET("/:id", controllers.GetSkuByID).
		POST("", controllers.CreateSku).
		DELETE("/:id", controllers.DeleteSku).
		PUT("/:id", controllers.UpdateSku)

	// Inventory routes
	server.Group("/inventories", middlewares.AuthMiddleware(false)).
		GET("", controllers.GetInventories).
		GET("/:id", controllers.GetInventoryByID).
		POST("", controllers.CreateInventory).
		DELETE("/:id", controllers.DeleteInventory).
		PUT("/:id", controllers.UpdateInventory)
}