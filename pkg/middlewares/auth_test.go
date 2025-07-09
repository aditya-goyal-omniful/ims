package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// dummyHandler is used to validate middleware pass-through
func dummyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		tenantIDHeader string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing X-Tenant-ID header",
			tenantIDHeader: "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing X-Tenant-ID header",
		},
		{
			name:           "invalid X-Tenant-ID UUID format",
			tenantIDHeader: "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid X-Tenant-ID format",
		},
		{
			name:           "valid X-Tenant-ID UUID format",
			tenantIDHeader: uuid.NewString(),
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.Use(AuthMiddleware())
			r.GET("/test", dummyHandler)

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.tenantIDHeader != "" {
				req.Header.Set("X-Tenant-ID", tt.tenantIDHeader)
			}
			c.Request = req

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
