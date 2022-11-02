package controllers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// EnvType is variable in .env file
type EnvType struct {
	DbConn   string
	HTTPPort string
}

// Env is global var for EnvType
var Env = EnvType{}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	Env.DbConn = getEnv("DB_CONN", "default conn")
	Env.HTTPPort = getEnv("HTTP_PORT", "3100")
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
