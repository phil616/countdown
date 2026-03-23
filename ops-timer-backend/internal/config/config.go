package config

import (
	"fmt"
	"os"
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

// Load 从配置文件与环境变量加载配置。环境变量前缀为 TIMER_，嵌套键中的点替换为下划线
//（例如 TIMER_AUTH_JWT_SECRET 对应 auth.jwt_secret）。
// 若 path 非空但文件不存在，则跳过文件，仅使用环境变量及下方默认值（适用于 Docker 等仅注入环境变量的场景）。
// path 为空时同样不读取文件，仅环境变量。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	v.SetEnvPrefix("TIMER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if path != "" {
		fi, err := os.Stat(path)
		switch {
		case err == nil && fi.IsDir():
			return nil, fmt.Errorf("config path %q is a directory", path)
		case err == nil && !fi.IsDir():
			v.SetConfigFile(path)
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("read config file %q: %w", path, err)
			}
		case err != nil && os.IsNotExist(err):
			// 无配置文件：仅 TIMER_* 环境变量 + 默认值
		default:
			return nil, fmt.Errorf("stat config path %q: %w", path, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
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
