package config

import (
	"log"
	"os"
	"truckx/internal/models"

	"github.com/joho/godotenv"
)

func Load() models.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file", err)
	}
	return models.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_DATABASE"),
		JWTKet:   os.Getenv("JWT_KEY"),
	}
}
