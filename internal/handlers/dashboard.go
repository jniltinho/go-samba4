package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *AppContext) DashboardGET(c *echo.Context) error {
	// Gather statistics from AD if connection is available
	users, err := app.LDAPClient.GetAllUsers("")
	userCount := len(users)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
	}

	groups, err := app.LDAPClient.GetAllGroups("")
	groupCount := len(groups)
	if err != nil {
		fmt.Printf("Error fetching groups: %v\n", err)
	}

	ous, err := app.LDAPClient.GetAllOUs()
	ouCount := len(ous)
	if err != nil {
		fmt.Printf("Error fetching OUs: %v\n", err)
	}

	return c.Render(http.StatusOK, "dashboard", map[string]interface{}{
		"TotalUsers":  userCount,
		"TotalGroups": groupCount,
		"TotalOUs":    ouCount,
	})
}
