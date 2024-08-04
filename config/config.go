package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

const (
	IS_DEVELOP_MODE = true
	APP_LOG_FOLDER  = "./logs/"
)

var (
	DatabaseName     = os.Getenv("DB_DATABASE")
	DatabasePassword = os.Getenv("DB_PASSWORD")
	DatabaseUsername = os.Getenv("DB_USERNAME")
	DatabasePort     = os.Getenv("DB_PORT")
	DatabaseHost     = os.Getenv("DB_HOST")
)
