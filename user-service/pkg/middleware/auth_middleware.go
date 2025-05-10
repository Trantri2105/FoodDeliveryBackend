package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-service/pkg/jwt"
)

type AuthMiddleware interface {
	ValidateAndExtractJwt() gin.HandlerFunc
}

const (
	JWTClaimsContextKey = "JWTClaimsContextKey"
)

type authMiddleware struct {
	jwt jwt.Utils
}

func (a *authMiddleware) ValidateAndExtractJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Authorization header is empty",
			})
			return
		}
		header := strings.Fields(authHeader)
		if len(header) != 2 && header[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Authorization header is invalid",
			})
			return
		}
		accessToken := header[1]
		claims, err := a.jwt.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}

		c.Set(JWTClaimsContextKey, claims)
		c.Next()
	}
}

func NewAuthMiddleware(jwtService jwt.Utils) AuthMiddleware {
	return &authMiddleware{jwt: jwtService}
}
