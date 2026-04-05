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

func Load(envFile ...string) {
	f := "../.env"
	if len(envFile) > 0 {
		f = envFile[0]
	}
	err := godotenv.Load(f)
	if err != nil {
		log.Fatal("error: loading .env:", err)
	}

	DBUrl = os.Getenv("DB_URL")
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
}
