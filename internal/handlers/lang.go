package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

var validLangs = map[string]bool{
	"en": true,
	"pt": true,
	"es": true,
}

// SetLanguage handles GET /lang/:code — sets the language cookie and redirects back.
func (app *AppContext) SetLanguage(c *echo.Context) error {
	code := c.Param("code")
	if !validLangs[code] {
		code = "en"
	}

	cookie := &http.Cookie{
		Name:     "samba4_lang",
		Value:    code,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	referer := c.Request().Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	return c.Redirect(http.StatusFound, referer)
}

// LangFromRequest detects the preferred language from cookie → Accept-Language → default "en".
func LangFromRequest(c *echo.Context) string {
	if cookie, err := c.Cookie("samba4_lang"); err == nil && validLangs[cookie.Value] {
		return cookie.Value
	}
	accept := c.Request().Header.Get("Accept-Language")
	if len(accept) >= 2 {
		prefix := accept[:2]
		switch prefix {
		case "pt":
			return "pt"
		case "es":
			return "es"
		}
	}
	return "en"
}

