package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	MySQL struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		Params   string
	}

	JWT struct {
		SecretKey    string
		TTLMinutes   int
		CookieDomain string
	}
}

func LoadConfig() (*Config, error) {
	var err error = godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load dotenv file: %w", err)
	}

	// Получение данных конфигурации приложения
	var config *Config = &Config{
		AppName: getEnvWithDefault("APP_NAME", "ToDoApp(develop)"),
		AppEnv:  getEnvWithDefault("APP_ENV", "dev"),
		AppPort: getEnvWithDefault("APP_PORT", "8080"),
	}

	// Получение данных конфигурации БД
	config.MySQL.Host = getEnvWithDefault("MYSQL_HOST", "127.0.0.1")
	config.MySQL.Port = getEnvWithDefault("MYSQL_PORT", "3306")
	config.MySQL.User = getEnvWithDefault("MYSQL_USER", "root")
	config.MySQL.Password = getEnvWithDefault("MYSQL_PASSWORD", "%")
	config.MySQL.Params = getEnvWithDefault("MYSQL_PARAMS", "charset=utf8mb4&parseTime=true&loc=Local")

	// Получение данных конфигурации JWT
	config.JWT.SecretKey = getEnvWithDefault("JWT_SECRET", "dev_jwt")
	ttl, err := strconv.Atoi(getEnvWithDefault("JWT_TTL_MINUTES", "120"))
	if err != nil {
		config.JWT.TTLMinutes = 120
	} else {
		config.JWT.TTLMinutes = ttl
	}
	config.JWT.CookieDomain = getEnvWithDefault("COOKIE_DOMAIN", "localhost")

	return config, nil
}

func (config *Config) GetMySQLDSN() string {
	// Формирование DSN строки для подключения к БД
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		config.MySQL.User, config.MySQL.Password,
		config.MySQL.Host, config.MySQL.Port,
		config.MySQL.DBName, config.MySQL.Params,
	)
}

func getEnvWithDefault(key string, defaultValue string) string {
	var value string = os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
