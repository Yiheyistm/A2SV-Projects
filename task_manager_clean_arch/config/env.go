package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration

type Env struct {
	AppEnv                 string
	ServerAddress          string
	ContextTimeout         int
	DBHost                 string
	DBHostURI              string
	DBPort                 string
	DBUserCollection       string
	DBTaskCollection       string
	DBPass                 string
	DBName                 string
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
}

// Load reads configuration from environment variables
func Load() *Env {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	env := &Env{
		AppEnv:                 GetEnvString("APP_ENV", "development"),
		ServerAddress:          GetEnvString("SERVER_ADDRESS", ":8080"),
		ContextTimeout:         GetEnvInt("CONTEXT_TIMEOUT", 30),
		DBHost:                 GetEnvString("DB_HOST", "localhost"),
		DBHostURI:              GetEnvString("DB_HOST_URI", "mongodb://localhost:27017"),
		DBPort:                 GetEnvString("DB_PORT", "27017"),
		DBUserCollection:       GetEnvString("DB_USER_COLLECTION", "users"),
		DBTaskCollection:       GetEnvString("DB_TASK_COLLECTION", "tasks"),
		DBPass:                 GetEnvString("DB_PASS", "password"),
		DBName:                 GetEnvString("DB_NAME", "task_manager"),
		AccessTokenExpiryHour:  GetEnvInt("ACCESS_TOKEN_EXPIRY_HOUR", 1),
		RefreshTokenExpiryHour: GetEnvInt("REFRESH_TOKEN_EXPIRY_HOUR", 24),
		AccessTokenSecret:      GetEnvString("ACCESS_TOKEN_SECRET", "secret"),
		RefreshTokenSecret:     GetEnvString("REFRESH_TOKEN_SECRET", "secret"),
	}

	return env
}

func GetEnvString(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
