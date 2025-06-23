package routes

import (
	"github.com/aditya-goyal-omniful/ims/pkg/controllers"
	"github.com/aditya-goyal-omniful/ims/pkg/middlewares"
	"github.com/omniful/go_commons/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(server *http.Server) {
	// Tenant Routes
	server.GET("tenants", controllers.GetTenants)
	server.GET("tenants/:id", controllers.GetTenantByID)
	server.POST("tenants", controllers.CreateTenant)
	server.DELETE("tenants/:id", controllers.DeleteTenant)
	server.PUT("tenants/:id", controllers.UpdateTenant)

	// Seller Routes
	server.GET("sellers", controllers.GetSellers)
	server.GET("sellers/:id", controllers.GetSellerByID)
	server.POST("sellers", controllers.CreateSeller)
	server.DELETE("sellers/:id", controllers.DeleteSeller)
	server.PUT("sellers/:id", controllers.UpdateSeller)

	// Hub routes
	server.Group("/hubs", middlewares.AuthMiddleware()).
		GET("", controllers.GetHubs).
		GET("/:id", controllers.GetHubByID).
		POST("", controllers.CreateHub).
		DELETE("/:id", controllers.DeleteHub).
		PUT("/:id", controllers.UpdateHub)

	// SKU routes (Tenant + Seller)
	server.Group("/skus", middlewares.AuthMiddleware()).
		GET("", controllers.GetSkus).
		GET("/:id", controllers.GetSkuByID).
		POST("", controllers.CreateSku).
		DELETE("/:id", controllers.DeleteSku).
		PUT("/:id", controllers.UpdateSku)

	// Inventory routes
	server.Group("/inventories", middlewares.AuthMiddleware()).
		GET("", controllers.GetInventories).
		GET("/:id", controllers.GetInventoryByID).
		POST("", controllers.CreateInventory).
		DELETE("/:id", controllers.DeleteInventory).
		PUT("/:id", controllers.UpdateInventory).
		POST("/upsert", controllers.UpsertInventory).
		GET("/view", controllers.ViewInventoryWithDefaults)


	// InterService Communication
	server.GET("validators/validate_order/:hub_id/:sku_id", controllers.ValidateOrder)
	server.POST("/inventory/check-and-update", controllers.CheckAndUpdateInventory)


	// Swagger Routes
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


}