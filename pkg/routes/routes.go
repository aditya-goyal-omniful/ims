package routes

import (
	"github.com/aditya-goyal-omniful/ims/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Hub routes
	r.GET("/hubs", controllers.GetHubs)
	r.GET("/hubs/:id", controllers.GetHubByID)
	r.POST("/hubs", controllers.CreateHub)
	r.DELETE("/hubs/:id", controllers.DeleteHub)
	r.PUT("/hubs/:id", controllers.UpdateHub)

	// SKU routes
	r.GET("/skus", controllers.GetSkus)
	r.GET("/skus/:id", controllers.GetSkuByID)
	r.POST("/skus", controllers.CreateSku)
	r.DELETE("/skus/:id", controllers.DeleteSku)
	r.PUT("/skus/:id", controllers.UpdateSku)

	// Inventory routes
	r.GET("/hubs", controllers.GetInventory)
	r.GET("/hubs/:id", controllers.GetInventoryByID)
	r.POST("/hubs", controllers.CreateInventory)
	r.DELETE("/hubs/:id", controllers.DeleteInventory)
	r.PUT("/hubs/:id", controllers.UpdateInventory)
}