package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/model"
)

// RequireRole returns middleware that restricts access to users whose role
// meets or exceeds the required minimum level.
//
//	Viewer        (1) — read-only dashboard, logs, status
//	Operator      (2) — config editing, sites, upstreams, certs, DNS, Lua
//	Administrator (3) — users, process control, settings, all Operator
func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := Claims(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		if !model.RequireRole(claims.Role, required) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}

// RequireAdmin is a shortcut for RequireRole(model.RoleAdministrator).
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(model.RoleAdministrator)
}

// RequireOperator is a shortcut for RequireRole(model.RoleOperator).
func RequireOperator() gin.HandlerFunc {
	return RequireRole(model.RoleOperator)
}
