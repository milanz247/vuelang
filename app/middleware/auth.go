package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth is a JWT Bearer token guard.
// Replace the stub with real JWT validation (e.g. golang-jwt/jwt).
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			return
		}
		// TODO: parse & verify the token, then set user context:
		// token := strings.TrimPrefix(header, "Bearer ")
		// claims, err := jwt.Verify(token, cfg.JWTSecret)
		// c.Set("user_id", claims.UserID)
		c.Next()
	}
}
