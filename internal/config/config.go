package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	LDAP     LDAPConfig     `mapstructure:"ldap"`
	Database DatabaseConfig `mapstructure:"database"`
	Session  SessionConfig  `mapstructure:"session"`
	Security SecurityConfig `mapstructure:"security"`
	RBAC     RBACConfig     `mapstructure:"rbac"`
}

type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
	DevMode bool   `mapstructure:"dev_mode"`
}

type LDAPConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	UseTLS        bool   `mapstructure:"use_tls"`
	SkipTLSVerify bool   `mapstructure:"skip_tls_verify"`
	BaseDN        string `mapstructure:"base_dn"`
	BindUser      string `mapstructure:"bind_user"`
	BindPass      string `mapstructure:"bind_pass"` // typically fed from env SAMBA4_LDAP_PASS
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Path   string `mapstructure:"path"`
	DSN    string `mapstructure:"dsn"`
}

type SessionConfig struct {
	Secret         string `mapstructure:"secret"`
	TimeoutMinutes int    `mapstructure:"timeout_minutes"`
	CookieSecure   bool   `mapstructure:"cookie_secure"`
	CookieSameSite string `mapstructure:"cookie_same_site"`
}

type SecurityConfig struct {
	MaxLoginAttempts int  `mapstructure:"max_login_attempts"`
	LockoutMinutes   int  `mapstructure:"lockout_minutes"`
	RequireTOTP      bool `mapstructure:"require_totp"`
}

type RBACConfig struct {
	AdminGroup    string `mapstructure:"admin_group"`
	OperatorGroup string `mapstructure:"operator_group"`
	HelpdeskGroup string `mapstructure:"helpdesk_group"`
	ReadonlyGroup string `mapstructure:"readonly_group"`
}

var globalConfig *Config

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/go-samba4/")
	}

	// Environment variable support
	v.SetEnvPrefix("SAMBA4")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind specific env vars
	_ = v.BindEnv("ldap.bind_pass", "SAMBA4_LDAP_PASS")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Set some defaults if not present
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Database.Driver == "" {
		cfg.Database.Driver = "sqlite"
		cfg.Database.Path = "data.db"
	}

	globalConfig = &cfg

	return &cfg, nil
}

func Get() *Config {
	return globalConfig
}
