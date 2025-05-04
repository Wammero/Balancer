package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis    RedisConfig  `yaml:"redis"`
	Server   ServerConfig `yaml:"server"`
	Backends []string     `yaml:"backends"`
	DB       DBConfig     `yaml:"db"`
	JWT      JWTConfig    `yaml:"jwt"`
}

type JWTConfig struct {
	SecretKey string
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	Timeout      time.Duration `yaml:"timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

func NewConfig() *Config {
	config := &Config{}
	err := loadConfigFromFile(config, "config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config from file: %v", err)
	}

	config.JWT.SecretKey = getEnv("JWT_SECRET_KEY", config.JWT.SecretKey)

	config.Redis.Host = getEnv("REDIS_HOST", config.Redis.Host)
	config.Redis.Port = getEnv("REDIS_PORT", config.Redis.Port)

	config.Server.Port = getEnv("SERVER_PORT", config.Server.Port)
	config.Server.Timeout = getEnvDuration("SERVER_TIMEOUT", config.Server.Timeout)
	config.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", config.Server.ReadTimeout)
	config.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", config.Server.WriteTimeout)
	config.Server.IdleTimeout = getEnvDuration("SERVER_IDLE_TIMEOUT", config.Server.IdleTimeout)

	config.DB.Host = getEnv("DB_HOST", config.DB.Host)
	config.DB.Port = getEnv("DB_PORT", config.DB.Port)
	config.DB.User = getEnv("DB_USER", config.DB.User)
	config.DB.Password = getEnv("DB_PASSWORD", config.DB.Password)
	config.DB.Name = getEnv("DB_NAME", config.DB.Name)

	backendsEnv := os.Getenv("BACKENDS")
	if backendsEnv != "" {
		config.Backends = parseBackends(backendsEnv)
	}

	return config
}

func (db DBConfig) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name)
}

func parseBackends(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func loadConfigFromFile(config *Config, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, config)
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Invalid duration value for %s: %v", key, err)
	}
	return duration
}
