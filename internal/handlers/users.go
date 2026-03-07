package handlers

import (
	"net/http"
	"go-samba4/internal/ldap"

	"github.com/labstack/echo/v5"
)

func (app *AppContext) UsersListGET(c *echo.Context) error {
	// Query params for pagination or filtering can go here
	users, err := app.LDAPClient.GetAllUsers("")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Filter out system accounts
	exclude := map[string]bool{
		"krbtgt":         true,
		"Guest":          true,
		"join-slave":     true,
		"join-backup":    true,
		"dns-ucs-adm-90": true,
	}

	var filteredUsers []ldap.User
	for _, u := range users {
		if !exclude[u.SAMAccountName] {
			filteredUsers = append(filteredUsers, u)
		}
	}

	return c.Render(http.StatusOK, "users/list", map[string]interface{}{
		"Users": filteredUsers,
	})
}

// Further implementations for Create, Edit, Detail
func (app *AppContext) UsersFormGET(c *echo.Context) error {
	return c.Render(http.StatusOK, "users/form", nil)
}

func (app *AppContext) UsersDetailGET(c *echo.Context) error {
	return c.Render(http.StatusOK, "users/detail", nil)
}
