package auth

import (
	"crypto/tls"
	"fmt"

	goldap "github.com/go-ldap/ldap/v3"
	"go-samba4/internal/config"
	"go-samba4/internal/ldap"
)

// AuthenticateUser attempts to bind to AD using the provided username and password
func AuthenticateUser(cfg *config.Config, username, password string) (*ldap.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	url := fmt.Sprintf("ldap://%s:%d", cfg.LDAP.Host, cfg.LDAP.Port)
	if cfg.LDAP.UseTLS && cfg.LDAP.Port == 636 {
		url = fmt.Sprintf("ldaps://%s:%d", cfg.LDAP.Host, cfg.LDAP.Port)
	}

	var conn *goldap.Conn
	var err error

	if cfg.LDAP.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.LDAP.SkipTLSVerify,
			ServerName:         cfg.LDAP.Host,
		}
		if cfg.LDAP.Port == 636 {
			conn, err = goldap.DialURL(url, goldap.DialWithTLSConfig(tlsConfig))
		} else {
			conn, err = goldap.DialURL(url)
			if err == nil {
				err = conn.StartTLS(tlsConfig)
			}
		}
	} else {
		conn, err = goldap.DialURL(url)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}
	defer conn.Close()

	// A more robust AD approach: bind with the service account, search for the user's DN, then bind with that DN + password.
	err = conn.Bind(cfg.LDAP.BindUser, cfg.LDAP.BindPass)
	if err != nil {
		return nil, fmt.Errorf("service account bind failed: %w", err)
	}

	// Search for the user to get their true DN
	searchRequest := goldap.NewSearchRequest(
		cfg.LDAP.BaseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", goldap.EscapeFilter(username)),
		[]string{"dn", ldap.AttrSAMAccountName, ldap.AttrDisplayName, ldap.AttrMemberOf},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil || len(sr.Entries) == 0 {
		return nil, fmt.Errorf("invalid credentials or user not found")
	}

	userDN := sr.Entries[0].DN

	// Now try binding with the user's actual DN and their password
	err = conn.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	u := &ldap.User{
		DN:             sr.Entries[0].DN,
		SAMAccountName: sr.Entries[0].GetAttributeValue(ldap.AttrSAMAccountName),
		DisplayName:    sr.Entries[0].GetAttributeValue(ldap.AttrDisplayName),
		MemberOf:       sr.Entries[0].GetAttributeValues(ldap.AttrMemberOf),
	}

	return u, nil
}
