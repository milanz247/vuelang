package middleware

import (
	"github.com/gin-gonic/gin"

	"vuelang/internal/framework/response"
)

// RequireRole restricts a route to users who have at least one of the given roles.
// The auth middleware must run first (sets "user_roles" in context).
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		userRoles, exists := c.Get("user_roles")
		if !exists {
			response.Forbidden(c)
			c.Abort()
			return
		}

		for _, r := range userRoles.([]string) {
			if allowed[r] {
				c.Next()
				return
			}
		}

		response.Forbidden(c)
		c.Abort()
	}
}
