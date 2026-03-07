package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-samba4/internal/config"
	"go-samba4/internal/models"
	// "go-samba4/internal/ldap"  // would be needed to check group mapping if not stored in session
)

// RBAC checks if the user has the required group role
func RBAC(requiredGroup string, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessionRaw := c.Get("session")
			if sessionRaw == nil {
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			// In a real AD scenario we would look up the user's groups via LDAP
			// For this skeleton, we assume `session.Username` implies some checks
			_ = sessionRaw.(*models.Session)

			// Role checking logic (stubbed out for now - would usually intersect session.MemberOf with requiredGroup)
			// Example: if !slices.Contains(user.MemberOf, requiredGroup) { return echo.ErrForbidden }

			return next(c)
		}
	}
}
