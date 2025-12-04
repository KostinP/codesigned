package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	Database   DBConfig         `yaml:"database"`
	Clickhouse ClickhouseConfig `yaml:"clickhouse"`
	Analytics  AnalyticsConfig  `yaml:"analytics"`
	Telegram   Telegram         `yaml:"telegram"`
	JWT        JWTConfig        `yaml:"jwt"`
	Container  ContainerConfig  `yaml:"container"`
	Logging    LoggingConfig    `yaml:"logging"`
	Cors       CorsConfig       `yaml:"cors"`
	Mode       string
}

type AppConfig struct {
	Env  string `yaml:"env"`
	Port int    `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

type ClickhouseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AnalyticsConfig struct {
	Enabled bool `yaml:"enabled"`
}

type Telegram struct {
	Token         string `yaml:"token"`
	WebhookDomain string `yaml:"webhook_domain"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type ContainerConfig struct {
	TimeoutSeconds int     `yaml:"timeout_seconds"`
	MemoryLimitMB  int     `yaml:"memory_limit_mb"`
	CPULimit       float64 `yaml:"cpu_limit"`
	BaseImage      string  `yaml:"base_image"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

func Load() *Config {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "dev"
	}

	path := "./configs/" + mode + ".yaml"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	// Расширение плейсхолдеров ${VAR} из env vars
	expanded := os.ExpandEnv(string(file))
	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		log.Fatalf("Не удалось распарсить конфигурационный файл: %v", err)
	}
	cfg.Mode = mode
	return &cfg
}
