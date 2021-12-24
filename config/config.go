package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	_model "github.com/fanajib5/url_shortener_PA/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DotEnv struct {
	DB_HOST       string
	DB_PORT       string
	DB_NAME       string
	DB_USER       string
	DB_PWD        string
	SECRET_HMAC   string
	SECRET_NANOID string
}

func (env *DotEnv) LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env.DB_HOST = os.Getenv("DB_HOST")
	env.DB_PORT = os.Getenv("DB_PORT")
	env.DB_NAME = os.Getenv("DB_NAME")
	env.DB_USER = os.Getenv("DB_USER")
	env.DB_PWD = os.Getenv("DB_PASSWORD")
	env.SECRET_HMAC = os.Getenv("SECRET_HMAC")
	env.SECRET_NANOID = os.Getenv("SECRET_NANOID_ALPHANUMERIC")
}

func InitDB() {
	env := &DotEnv{}
	env.LoadDotEnv()

	var errConDB error

	dsn := env.DB_USER + ":" + env.DB_PWD + "@tcp(" + env.DB_HOST + ":" + env.DB_PORT + ")/" + env.DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"

	DB, errConDB = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errConDB != nil {
		panic(errConDB.Error())
	}
	initMigrate()
}

func initMigrate() {
	DB.AutoMigrate(&_model.User{})
	DB.AutoMigrate(&_model.UrlData{})
}
