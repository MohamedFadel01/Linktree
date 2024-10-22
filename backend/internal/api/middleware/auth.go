package middleware

import (
	"linktree-mohamedfadel-backend/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateJWTFromContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}

func OptionalJWTFromContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateJWT(tokenString)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				ctx.Abort()
				return
			}

			ctx.Set("username", claims.Username)
			ctx.Next()
		}

		ctx.Set("username", "anonymous")
		ctx.Next()
	}
}
