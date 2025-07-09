package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func dummyHandler2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged"})
}

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		tenantIDHeader string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "request logging with tenant header",
			tenantIDHeader: "123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
			expectedBody:   "logged",
		},
		{
			name:           "request logging without tenant header",
			tenantIDHeader: "",
			expectedStatus: http.StatusOK,
			expectedBody:   "logged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(RequestLogger())
			r.GET("/ping", dummyHandler2)

			req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
			if tt.tenantIDHeader != "" {
				req.Header.Set("X-Tenant-ID", tt.tenantIDHeader)
			}

			start := time.Now()
			r.ServeHTTP(w, req)
			_ = time.Since(start) // just to simulate time-based logger execution

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
