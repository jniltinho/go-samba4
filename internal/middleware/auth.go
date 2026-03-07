package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-samba4/internal/auth"
)

func RequireAuth(sm *auth.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SessionCookieName)
			if err != nil {
				// No session cookie found
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			session, err := sm.GetSession(cookie.Value)
			if err != nil || session == nil {
				// Invalid or expired session
				sm.ClearCookie(c.Response().Writer)
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			// Store session and username in context
			c.Set("session", session)
			c.Set("username", session.Username)

			return next(c)
		}
	}
}

// OptionalAuth attaches session if available, but doesn't block
func OptionalAuth(sm *auth.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(auth.SessionCookieName)
			if err == nil {
				session, err := sm.GetSession(cookie.Value)
				if err == nil && session != nil {
					c.Set("session", session)
					c.Set("username", session.Username)
				}
			}
			return next(c)
		}
	}
}
