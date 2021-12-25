package config

import (
	"fmt"

	_model "github.com/fanajib5/url_shortener_PA/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	dbName := "url_shortener"

	config := map[string]string{
		"DB_Username": "root",
		"DB_Password": "",
		"DB_Port":     "3306",
		"DB_Host":     "127.0.0.1",
		"DB_Name":     dbName,
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config["DB_Username"],
		config["DB_Password"],
		config["DB_Host"],
		config["DB_Port"],
		config["DB_Name"])

	var errConDB error
	DB, errConDB = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if errConDB != nil {
		panic(errConDB)
	}
	initMigrate()
}

func initMigrate() {
	DB.AutoMigrate(&_model.User{})
	DB.AutoMigrate(&_model.UrlData{})
}
