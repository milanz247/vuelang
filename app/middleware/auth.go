package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	jwtpkg "vuelang/internal/framework/jwt"
	"vuelang/internal/framework/response"
)

// AuthMiddleware validates JWT Bearer tokens and injects user context.
type AuthMiddleware struct {
	jwt jwtpkg.Service
}

func NewAuth(jwt jwtpkg.Service) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

// Handle returns the Gin middleware function.
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.jwt.ValidateAccess(tokenStr)
		if err != nil {
			if err == jwtpkg.ErrTokenExpired {
				response.Unauthorized(c, "Access token has expired")
			} else {
				response.Unauthorized(c, "Invalid access token")
			}
			c.Abort()
			return
		}

		// Inject verified claims into context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Next()
	}
}
