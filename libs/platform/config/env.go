package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	return godotenv.Load(path)
}

func GetString(key string) string {
	return os.Getenv(key)
}
