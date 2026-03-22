package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	DB   DBConfig
	JWT  JWTConfig
	Seed SeedConfig
}

type AppConfig struct {
	Name         string
	Env          string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	AutoMigrate  bool
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	Secret         string
	ExpiresInHours int
}

type SeedConfig struct {
	AdminEmail       string
	AdminPassword    string
	AdminDisplayName string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:         getEnv("APP_NAME", "devfolio-api"),
			Env:          getEnv("APP_ENV", "local"),
			Port:         getEnv("APP_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
			AutoMigrate:  getEnvAsBool("AUTO_MIGRATE", true),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "devfolio"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Bangkok"),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "super-secret-change-me"),
			ExpiresInHours: getEnvAsInt("JWT_EXPIRES_IN_HOURS", 24),
		},
		Seed: SeedConfig{
			AdminEmail:       getEnv("ADMIN_EMAIL", "admin@example.com"),
			AdminPassword:    getEnv("ADMIN_PASSWORD", "changeme123"),
			AdminDisplayName: getEnv("ADMIN_DISPLAY_NAME", "Admin"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("invalid int for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("invalid bool for %s, using default %v", key, defaultValue)
		return defaultValue
	}
	return value
}
