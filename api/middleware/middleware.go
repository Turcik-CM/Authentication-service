package middleware

import (
	"auth-service/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{"https://turk.gophers.com", "https://turk.gophers.com"}
		origin := c.Request.Header.Get("Origin")

		// Проверяем, разрешен ли домен
		for _, o := range allowedOrigins {
			if origin == o {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Остальные CORS заголовки
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Обработка OPTIONS запросов (preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Передаем управление дальше
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
