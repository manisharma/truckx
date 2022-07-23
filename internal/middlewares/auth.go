package middlewares

import (
	"strings"
	"truckx/internal/models"
	"truckx/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func Auth(jwtKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "request does not contain an access token"})
			ctx.Abort()
			return
		}
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 {
			ctx.JSON(401, gin.H{"error": "access token is not valid"})
			ctx.Abort()
			return
		}
		user, err := auth.GetUserFromToken(tokenParts[1], []byte(jwtKey))
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set(string(models.UserCtxKey), user)
		ctx.Next()
	}
}
