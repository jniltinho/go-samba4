package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *AppContext) UsersListGET(c echo.Context) error {
	// Query params for pagination or filtering can go here
	users, err := app.LDAPClient.GetAllUsers("")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "users/list", map[string]interface{}{
		"Users": users,
	})
}

// Further implementations for Create, Edit, Detail
func (app *AppContext) UsersFormGET(c echo.Context) error {
	return c.Render(http.StatusOK, "users/form", nil)
}

func (app *AppContext) UsersDetailGET(c echo.Context) error {
	return c.Render(http.StatusOK, "users/detail", nil)
}
