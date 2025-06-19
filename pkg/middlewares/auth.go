package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")

		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Tenant-ID header"})
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(tenantID); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid X-Tenant-ID format"})
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}
