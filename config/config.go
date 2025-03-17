package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSigninMethod     *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
	CorsConfig
}

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host:     "localhost",
		Port:     "3306",
		Database: "go_api_backend",
		Username: "root",
		Password: "",
		Driver:   "mysql",
	}

	c.APIConfig = APIConfig{
		ApiPort: "8888",
	}

	accessTokenLifeTime := time.Duration(1) * time.Hour

	c.TokenConfig = TokenConfig{
		ApplicationName:     "Go API Backend",
		JwtSignatureKey:     []byte("random-secret-key"),
		JwtSigninMethod:     jwt.SigningMethodHS256,
		AccessTokenLifeTime: accessTokenLifeTime,
	}

	c.CorsConfig = CorsConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
