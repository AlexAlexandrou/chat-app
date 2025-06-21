package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func GetEnv(envVar string, defaultValue string) string {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		return defaultValue
	}
	return value
}
