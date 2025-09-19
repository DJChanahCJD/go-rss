package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
		DBUrl: os.Getenv("DB_URL"),
		Port:  os.Getenv("PORT"),
	}
}
