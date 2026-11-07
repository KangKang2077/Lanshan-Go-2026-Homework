package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"week_06_web/auth"
)

const (
	RequestIDKey = "request_id"
	UserIDKey    = "user_id"
)

func NewRequestID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "Unkown"
	}
	return hex.EncodeToString(buf[:])
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = NewRequestID()
		}

		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.Info("http request",
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"request_id", c.GetString(RequestIDKey))
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func JWTAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			abortUnauthorized(c)
			return
		}

		claims, err := tokens.Parse(tokenString)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if token, err := c.Cookie(auth.CookieName); err == nil && token != "" {
		return token
	}

	header := c.GetHeader("Authorization")
	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code": "UNAUTHORIZED",
		"msg":  "unauthorized",
	})
}
