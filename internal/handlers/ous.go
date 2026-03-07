package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *AppContext) OUsTreeGET(c *echo.Context) error {
	ous, err := app.LDAPClient.GetAllOUs()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "ous/tree", map[string]interface{}{
		"OUs": ous,
	})
}
