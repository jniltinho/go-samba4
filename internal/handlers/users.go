package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"go-samba4/internal/ldap"
	"go-samba4/internal/models"

	"github.com/labstack/echo/v5"
)

// systemAccounts holds AD built-in accounts to exclude from the list.
var systemAccounts = map[string]bool{
	"krbtgt":         true,
	"Guest":          true,
	"join-slave":     true,
	"join-backup":    true,
	"dns-ucs-adm-90": true,
}

// writeAudit persists an audit log entry to the database.
func (app *AppContext) writeAudit(c *echo.Context, action, objectDN, details string) {
	username, _ := c.Get("username").(string)
	ip := c.RealIP()
	entry := models.AuditLog{
		AdminUser: username,
		IPAddress: ip,
		Action:    action,
		ObjectDN:  objectDN,
		Details:   details,
	}
	if err := app.DB.Create(&entry).Error; err != nil {
		slog.Warn("Failed to write audit log", "error", err)
	}
}

// UsersListGET renders the paginated user list.
func (app *AppContext) UsersListGET(c *echo.Context) error {
	users, err := app.LDAPClient.GetAllUsers("")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var filteredUsers []ldap.User
	for _, u := range users {
		if !systemAccounts[u.SAMAccountName] {
			filteredUsers = append(filteredUsers, u)
		}
	}

	isAdmin, _ := c.Get("is_admin").(bool)
	currentUser, _ := c.Get("username").(string)
	return c.Render(http.StatusOK, "users/list", map[string]interface{}{
		"Users":       filteredUsers,
		"IsAdmin":     isAdmin,
		"CurrentUser": currentUser,
	})
}

// UsersFormGET renders the create or edit form.
// When an :id param is present it pre-populates the form with the existing user.
func (app *AppContext) UsersFormGET(c *echo.Context) error {
	sam := c.Param("id")
	if sam == "" {
		// New user form — default: account enabled, pwd expires
		ous, _ := app.LDAPClient.GetAllOUs()
		return c.Render(http.StatusOK, "users/form", map[string]interface{}{
			"OUs":            ous,
			"AccountEnabled": true,
			"Domain":         app.LDAPClient.Domain(),
		})
	}

	// Edit form — load the user
	user, err := app.LDAPClient.GetUserBySAM(sam)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
	}

	isAdmin, _ := c.Get("is_admin").(bool)
	ous, _ := app.LDAPClient.GetAllOUs()
	return c.Render(http.StatusOK, "users/form", map[string]interface{}{
		"User":            user,
		"OUs":             ous,
		"IsAdmin":         isAdmin,
		"AccountEnabled":  (user.UserAccountControl & 2) == 0,
		"PwdNeverExpires": (user.UserAccountControl & 65536) != 0,
	})
}

// UsersDetailGET renders the user detail page.
func (app *AppContext) UsersDetailGET(c *echo.Context) error {
	sam := c.Param("id")
	user, err := app.LDAPClient.GetUserBySAM(sam)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
	}
	isAdmin, _ := c.Get("is_admin").(bool)
	currentUser, _ := c.Get("username").(string)
	return c.Render(http.StatusOK, "users/detail", map[string]interface{}{
		"User":        user,
		"IsAdmin":     isAdmin,
		"CurrentUser": currentUser,
	})
}

// UsersCreatePOST creates a new AD user from the submitted form.
func (app *AppContext) UsersCreatePOST(c *echo.Context) error {
	sam := c.FormValue("SAMAccountName")
	password := c.FormValue("Password")
	ouDN := c.FormValue("OUDN")

	if sam == "" || password == "" || ouDN == "" {
		ous, _ := app.LDAPClient.GetAllOUs()
		return c.Render(http.StatusUnprocessableEntity, "users/form", map[string]interface{}{
			"OUs":   ous,
			"Error": "Username, password, and OU are required",
		})
	}

	// Compute desired UAC from checkboxes (same logic as update).
	newUAC := 512
	if c.FormValue("AccountEnabled") != "1" {
		newUAC |= 2
	}
	if c.FormValue("PwdNeverExpires") == "1" {
		newUAC |= 65536
	}

	u := ldap.User{
		SAMAccountName:     sam,
		DisplayName:        c.FormValue("DisplayName"),
		GivenName:          c.FormValue("GivenName"),
		SN:                 c.FormValue("SN"),
		UserPrincipalName:  c.FormValue("UserPrincipalName"),
		Mail:               c.FormValue("Mail"),
		TelephoneNumber:    c.FormValue("TelephoneNumber"),
		Title:              c.FormValue("Title"),
		Department:         c.FormValue("Department"),
		UserAccountControl: newUAC,
	}

	if err := app.LDAPClient.CreateUser(u, password, ouDN); err != nil {
		slog.Error("Failed to create user", "sam", sam, "error", err)
		ous, _ := app.LDAPClient.GetAllOUs()
		return c.Render(http.StatusUnprocessableEntity, "users/form", map[string]interface{}{
			"User":  u,
			"OUs":   ous,
			"Error": fmt.Sprintf("Failed to create user: %v", err),
		})
	}

	displayName := u.DisplayName
	if displayName == "" {
		displayName = sam
	}
	app.writeAudit(c, "CREATE_USER",
		fmt.Sprintf("CN=%s,%s", displayName, ouDN),
		fmt.Sprintf("sAMAccountName=%s", sam),
	)

	return c.Redirect(http.StatusFound, "/users/"+sam)
}

// UsersUpdatePOST updates an existing AD user's attributes.
func (app *AppContext) UsersUpdatePOST(c *echo.Context) error {
	sam := c.Param("id")

	user, err := app.LDAPClient.GetUserBySAM(sam)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
	}

	// Compute UserAccountControl from checkboxes.
	// Base: 512 (NORMAL_ACCOUNT). Bit 1 = ACCOUNTDISABLE. Bit 16 = DONT_EXPIRE_PASSWORD.
	newUAC := 512
	if c.FormValue("AccountEnabled") != "1" {
		newUAC |= 2 // set ACCOUNTDISABLE bit
	}
	if c.FormValue("PwdNeverExpires") == "1" {
		newUAC |= 65536 // set DONT_EXPIRE_PASSWORD bit
	}

	updated := ldap.User{
		DisplayName:        c.FormValue("DisplayName"),
		GivenName:          c.FormValue("GivenName"),
		SN:                 c.FormValue("SN"),
		UserPrincipalName:  c.FormValue("UserPrincipalName"),
		Mail:               c.FormValue("Mail"),
		TelephoneNumber:    c.FormValue("TelephoneNumber"),
		Title:              c.FormValue("Title"),
		Department:         c.FormValue("Department"),
		UserAccountControl: newUAC,
	}

	if err := app.LDAPClient.UpdateUser(user.DN, updated); err != nil {
		slog.Error("Failed to update user", "sam", sam, "error", err)
		ous, _ := app.LDAPClient.GetAllOUs()
		return c.Render(http.StatusUnprocessableEntity, "users/form", map[string]interface{}{
			"User":  user,
			"OUs":   ous,
			"Error": fmt.Sprintf("Failed to update user: %v", err),
		})
	}

	// Optional password change — never logged
	if newPass := c.FormValue("Password"); newPass != "" {
		if err := app.LDAPClient.SetPassword(user.DN, newPass); err != nil {
			slog.Error("Failed to set password", "sam", sam, "error", err)
		}
	}

	app.writeAudit(c, "UPDATE_USER", user.DN,
		fmt.Sprintf("Updated attributes for sAMAccountName=%s", sam),
	)

	return c.Redirect(http.StatusFound, "/users/"+sam)
}

// UsersDeletePOST removes an AD user after confirmation of the sAMAccountName.
func (app *AppContext) UsersDeletePOST(c *echo.Context) error {
	sam := c.Param("id")
	confirmation := c.FormValue("confirm_name")

	if confirmation != sam {
		return c.String(http.StatusBadRequest, "Confirmation does not match username")
	}

	user, err := app.LDAPClient.GetUserBySAM(sam)
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
	}

	dn := user.DN
	if err := app.LDAPClient.DeleteUser(dn); err != nil {
		slog.Error("Failed to delete user", "sam", sam, "error", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to delete user: %v", err))
	}

	app.writeAudit(c, "DELETE_USER", dn,
		fmt.Sprintf("Deleted sAMAccountName=%s", sam),
	)

	return c.Redirect(http.StatusFound, "/users")
}
