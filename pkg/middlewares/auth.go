package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextTenantIDKey = "tenant_id"
	ContextSellerIDKey = "seller_id"
)

func AuthMiddleware(requireSeller bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		sellerID := c.GetHeader("X-Seller-ID")

		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Tenant-ID header"})
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(tenantID); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid X-Tenant-ID format"})
			return
		}

		if requireSeller {
			if sellerID == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Seller-ID header"})
				return
			}
			if _, err := uuid.Parse(sellerID); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid X-Seller-ID format"})
				return
			}
			c.Set("seller_id", sellerID)
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}
