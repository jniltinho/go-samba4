package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v5"
)

// CSRF middleware generates a token on GET and validates it on POST/PUT/DELETE
func CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()

			if req.Method == http.MethodGet {
				// Generate new token or use existing one from context if already set
				tokenBytes := make([]byte, 32)
				rand.Read(tokenBytes)
				token := base64.URLEncoding.EncodeToString(tokenBytes)
				c.Set("csrf", token)

				// Set token in a secondary cookie to be validated later
				http.SetCookie(c.Response(), &http.Cookie{
					Name:     "samba4_csrf",
					Value:    token,
					Path:     "/",
					HttpOnly: true,
					Secure:   req.TLS != nil,
					SameSite: http.SameSiteStrictMode,
				})

			} else if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
				// Check token
				cookie, err := c.Cookie("samba4_csrf")
				if err != nil {
					return echo.NewHTTPError(http.StatusForbidden, "CSRF cookie missing")
				}

				formToken := c.FormValue("_csrf")
				headerToken := req.Header.Get("X-CSRF-Token")

				providedToken := formToken
				if providedToken == "" {
					providedToken = headerToken
				}

				if providedToken == "" || providedToken != cookie.Value {
					return echo.NewHTTPError(http.StatusForbidden, "invalid CSRF token")
				}
			}

			return next(c)
		}
	}
}
