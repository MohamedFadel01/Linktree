package tests

import (
	"linktree-mohamedfadel-backend/internal/api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateJWTFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header missing",
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1vaGFtZWRmYWRlbCIsImV4cCI6MTcyOTYyNzczN30.qt-XB0XOHpFdR-yQxEXeNwl-P-QPccScYa2iy9t6Jhs",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid or expired token",
		},
		{
			name:           "malformed token",
			authHeader:     "Bearer malformed",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid or expired token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest("GET", "/", nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			c.Request = req

			middleware := middleware.ValidateJWTFromContext()
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestOptionalJWTFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		authHeader      string
		expectedStatus  int
		expectedError   string
		expectedContext string
	}{
		{
			name:            "no auth header",
			authHeader:      "",
			expectedStatus:  http.StatusOK,
			expectedContext: "anonymous",
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1vaGFtZWRmYWRlbCIsImV4cCI6MTcyOTYyNzczN30.qt-XB0XOHpFdR-yQxEXeNwl-P-QPccScYa2iy9t6Jhs",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid or expired token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest("GET", "/", nil)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			c.Request = req

			middleware := middleware.OptionalJWTFromContext()
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			if tt.expectedContext != "" {
				username, exists := c.Get("username")
				assert.True(t, exists)
				assert.Equal(t, tt.expectedContext, username)
			}
		})
	}
}
