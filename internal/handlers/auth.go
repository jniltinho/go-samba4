package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"samba4-admin/internal/auth"
)

// AuthLoginGET handles showing the login page
func (app *AppContext) AuthLoginGET(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/login", nil)
}

// AuthLoginPOST handles processing the login form
func (app *AppContext) AuthLoginPOST(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Try to authenticate via LDAP
	u, err := auth.AuthenticateUser(app.Config, username, password)
	if err != nil {
		return c.Render(http.StatusBadRequest, "auth/login", map[string]interface{}{
			"Error": "Invalid credentials or LDAP error",
		})
	}

	// Create session
	session, err := app.SessionMgr.CreateSession(u.SAMAccountName, c.RealIP(), c.Request().UserAgent())
	if err != nil {
		return c.Render(http.StatusInternalServerError, "auth/login", map[string]interface{}{
			"Error": "Failed to create session",
		})
	}

	app.SessionMgr.SetCookie(c.Response().Writer, session.Token)

	return c.Redirect(http.StatusFound, "/dashboard")
}

// AuthLogout handles ending the session
func (app *AppContext) AuthLogout(c echo.Context) error {
	sessionRaw := c.Get("session")
	if sessionRaw != nil {
		// remove from DB
		// app.SessionMgr.DeleteSession(sessionRaw.(*models.Session).Token)
	}
	app.SessionMgr.ClearCookie(c.Response().Writer)
	return c.Redirect(http.StatusFound, "/auth/login")
}
