package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	DATABASE_USER = os.Getenv("DATABASE_USER")
	DATABASE_PASS = os.Getenv("DATABASE_PASS")
	DATABASE_HOST = os.Getenv("DATABASE_HOST")
	DATABASE_PORT = os.Getenv("DATABASE_PORT")
	DATABASE      = os.Getenv("DATABASE")
)

var (
	JWT_ISSUER string = os.Getenv("JWT_ISSUER")
	JWT_SECRET string = os.Getenv("JWT_SECRET")
)

var (
	PORT string = os.Getenv("PORT")
)
