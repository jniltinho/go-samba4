package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *AppContext) GroupsListGET(c *echo.Context) error {
	groups, err := app.LDAPClient.GetAllGroups("")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "groups/list", map[string]interface{}{
		"Groups": groups,
	})
}

func (app *AppContext) GroupsFormGET(c *echo.Context) error {
	return c.Render(http.StatusOK, "groups/form", nil)
}
