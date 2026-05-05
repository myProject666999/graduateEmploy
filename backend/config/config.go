package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	ServerMode     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	UploadPath     string
}

var AppConfig *Config

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
	}

	AppConfig = &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ServerMode:     getEnv("SERVER_MODE", "debug"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", "123456"),
		DBName:         getEnv("DB_NAME", "graduate_employ"),
		JWTSecret:      getEnv("JWT_SECRET", "default_secret_key"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		UploadPath:     getEnv("UPLOAD_PATH", "./uploads"),
	}

	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
