package middleware

import (
	"auth-service/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
