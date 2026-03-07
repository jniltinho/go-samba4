package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"samba4-admin/internal/models"
)

func (app *AppContext) SettingsGET(c echo.Context) error {
	var settings []models.Setting
	app.DB.Find(&settings)

	return c.Render(http.StatusOK, "settings", map[string]interface{}{
		"Settings": settings,
	})
}

func (app *AppContext) SettingsPOST(c echo.Context) error {
	// Parse settings update form
	return c.Redirect(http.StatusFound, "/settings")
}
