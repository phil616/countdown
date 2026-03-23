package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
	OAuth     OAuthConfig     `mapstructure:"oauth"`
	SMTP      SMTPConfig      `mapstructure:"smtp"`
	Log       LogConfig       `mapstructure:"log"`
}

type ServerConfig struct {
	Host        string   `mapstructure:"host"`
	Port        int      `mapstructure:"port"`
	CorsOrigins []string `mapstructure:"cors_origins"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type AuthConfig struct {
	JWTSecret         string `mapstructure:"jwt_secret"`
	JWTExpiryHours    int    `mapstructure:"jwt_expiry_hours"`
	LoginLockAttempts int    `mapstructure:"login_lock_attempts"`
	LoginLockMinutes  int    `mapstructure:"login_lock_minutes"`
}

type SchedulerConfig struct {
	NotificationScanInterval string `mapstructure:"notification_scan_interval"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// Enabled 返回是否启用 SMTP 邮件通知
func (c *SMTPConfig) Enabled() bool {
	return c.Host != "" && c.From != ""
}

type OAuthConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	IssuerURL    string   `mapstructure:"issuer_url"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	FrontendURL  string   `mapstructure:"frontend_url"`
	Scopes       []string `mapstructure:"scopes"`
	AdminEmails  []string `mapstructure:"admin_emails"`
}

// IsConfigured 判断 OAuth 是否已完整配置
func (c *OAuthConfig) IsConfigured() bool {
	return c.Enabled && c.IssuerURL != "" && c.ClientID != "" && c.ClientSecret != ""
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("TIMER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Auth.JWTExpiryHours == 0 {
		cfg.Auth.JWTExpiryHours = 24
	}
	if cfg.Auth.LoginLockAttempts == 0 {
		cfg.Auth.LoginLockAttempts = 5
	}
	if cfg.Auth.LoginLockMinutes == 0 {
		cfg.Auth.LoginLockMinutes = 15
	}
	if cfg.Scheduler.NotificationScanInterval == "" {
		cfg.Scheduler.NotificationScanInterval = "10m"
	}
	if cfg.SMTP.Port == 0 {
		cfg.SMTP.Port = 587
	}

	return &cfg, nil
}
