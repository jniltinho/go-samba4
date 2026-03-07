package ldap

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"time"

	goldap "github.com/go-ldap/ldap/v3"
	"go-samba4/internal/config"
)

// Client wraps the LDAP connection
type Client struct {
	conn   *goldap.Conn
	config *config.LDAPConfig
}

// NewClient establishes a new bound connection to the LDAP server
func NewClient(cfg *config.LDAPConfig) (*Client, error) {
	url := fmt.Sprintf("ldap://%s:%d", cfg.Host, cfg.Port)
	if cfg.UseTLS && cfg.Port == 636 {
		url = fmt.Sprintf("ldaps://%s:%d", cfg.Host, cfg.Port)
	}

	dialOpts := []goldap.DialOpt{
		goldap.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}),
	}

	var conn *goldap.Conn
	var err error

	if cfg.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.SkipTLSVerify,
			ServerName:         cfg.Host,
		}
		if cfg.Port == 636 {
			conn, err = goldap.DialURL(url, goldap.DialWithTLSConfig(tlsConfig))
		} else {
			conn, err = goldap.DialURL(url, dialOpts...)
			if err == nil {
				err = conn.StartTLS(tlsConfig)
			}
		}
	} else {
		conn, err = goldap.DialURL(url, dialOpts...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}

	// Bind
	err = conn.Bind(cfg.BindUser, cfg.BindPass)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to bind to LDAP: %w", err)
	}

	return &Client{conn: conn, config: cfg}, nil
}

// Close closes the underlying LDAP connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Reconnect attempts to silently reconnect and rebind
func (c *Client) Reconnect() error {
	c.Close()

	newClient, err := NewClient(c.config)
	if err != nil {
		return err
	}

	c.conn = newClient.conn
	slog.Info("Successfully reconnected to LDAP")
	return nil
}

// Search with auto-reconnect fallback
func (c *Client) Search(searchRequest *goldap.SearchRequest) (*goldap.SearchResult, error) {
	res, err := c.conn.Search(searchRequest)
	if err != nil {
		if goldap.IsErrorWithCode(err, goldap.ErrorNetwork) {
			slog.Warn("LDAP connection lost during search, attempting reconnect...")
			if reconnErr := c.Reconnect(); reconnErr == nil {
				return c.conn.Search(searchRequest)
			}
		}
		return res, err
	}
	return res, nil
}
