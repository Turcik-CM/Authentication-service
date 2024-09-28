package middleware

import (
	"auth-service/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from a specific origin (e.g., https://example.com)
		// or use "*" for all origins (which isn't secure for credentials)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://your-allowed-domain.com")

		// Security headers for HTTPS
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}

func GetAccessTokenMid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		claims, err := token.ExtractClaimsRefresh(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)

		ctx.Next()
	}
}
