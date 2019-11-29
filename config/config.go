package config

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

//DBConfig ...
var DBConfig DatabaseConfig

//LoadConfig ...
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	initilizeDBConfig()
}

func initilizeDBConfig() {
	DBConfig = DatabaseConfig{
		Host:       os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_DATABASE"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}
