package config

import (
	"os"

	_model "github.com/fanajib5/url_shortener_PA/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var errConDB error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := "root:@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
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
