package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")

		if tenantID == "" {
			c.AbortWithStatusJSON(int(http.StatusUnauthorized), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Missing X-Tenant-ID header")})
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(tenantID); err != nil {
			c.AbortWithStatusJSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid X-Tenant-ID format")})
			return
		}

		c.Set("tenant_id", tenantID)
		c.Next()
	}
}
