package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUrl     string
	JWTSecret []byte
)

func Load() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("error: loading .env:", err)
	}

	DBUrl = os.Getenv("DB_URL")
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
}
