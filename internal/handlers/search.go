package handlers

import (
	"fmt"
	"net/http"

	goldap "github.com/go-ldap/ldap/v3"
	"github.com/labstack/echo/v5"
)

func (app *AppContext) SearchGET(c *echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.Render(http.StatusOK, "search", map[string]interface{}{
			"Results": []interface{}{},
		})
	}

	// Simple wildcard search across common attributes
	escapedQuery := goldap.EscapeFilter(fmt.Sprintf("*%s*", query))
	filter := fmt.Sprintf("(|(sAMAccountName=%s)(displayName=%s)(mail=%s))", escapedQuery, escapedQuery, escapedQuery)

	users, _ := app.LDAPClient.GetAllUsers(filter)
	groups, _ := app.LDAPClient.GetAllGroups(filter)

	return c.Render(http.StatusOK, "search", map[string]interface{}{
		"Query":  query,
		"Users":  users,
		"Groups": groups,
	})
}
