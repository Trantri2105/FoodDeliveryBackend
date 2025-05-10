package middlewares

import (
	"net/http"
	"runtime/debug"
	"strings"
	"thanhnt208/delivery-service/pkg/jwt"
	"thanhnt208/delivery-service/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// -------- Logging Middleware --------
func LoggingMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Ensure request_id exists
		requestID := c.Writer.Header().Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Writer.Header().Set("X-Request-ID", requestID)
		}

		c.Next()

		duration := time.Since(start)
		logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"user_agent", c.Request.UserAgent(),
			"remote_addr", c.ClientIP(),
			"request_id", requestID,
		)
	}
}

// -------- CORS Middleware --------
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // nếu cần

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// -------- Recovery Middleware --------
func RecoveryMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recover from panic",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(debug.Stack()),
				)

				// Đảm bảo header đúng khi trả lỗi
				c.Writer.Header().Set("Content-Type", "application/json")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}()

		c.Next()
	}
}

// -------- JWT Auth Middleware --------

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
		if len(header) != 2 || header[0] != "Bearer" {
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
			return
		}

		c.Set(JWTClaimsContextKey, claims)
		c.Next()
	}
}

func NewAuthMiddleware(jwtService jwt.Utils) AuthMiddleware {
	return &authMiddleware{jwt: jwtService}
}