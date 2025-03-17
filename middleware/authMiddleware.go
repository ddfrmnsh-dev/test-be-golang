package middleware

import (
	"api-backend-go/model"
	"api-backend-go/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(allowedRoles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}

func (a *authMiddleware) RequireToken(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var aH authHeader

		if err := ctx.ShouldBindHeader(&aH); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "status": false})
			return
		}

		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)

		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "status": false})
			return
		}

		newId, _ := strconv.Atoi(tokenClaim.UserId)

		validRole := false

		for _, role := range allowedRoles {
			if role == tokenClaim.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Forbidden Resource", "status": false})
			return
		}

		ctx.Set("user", model.User{Id: newId, Role: tokenClaim.Role})
		ctx.Next()
	}
}
