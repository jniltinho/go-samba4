package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go-samba4/internal/models"
)

// RequireAdmin blocks access for non-admin users (403 Forbidden).
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			session, ok := c.Get("session").(*models.Session)
			if !ok || session == nil {
				return c.Redirect(http.StatusFound, "/auth/login")
			}
			if !session.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "Administrator access required")
			}
			return next(c)
		}
	}
}

// RequireAdminOrSelf allows admins to access any resource, and normal users
// to access only their own resource (matched via the ":id" path parameter).
func RequireAdminOrSelf() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			session, ok := c.Get("session").(*models.Session)
			if !ok || session == nil {
				return c.Redirect(http.StatusFound, "/auth/login")
			}
			if session.IsAdmin || c.Param("id") == session.Username {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}
	}
}
