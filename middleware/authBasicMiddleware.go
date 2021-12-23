package middleware

import (
	_config "github.com/fanajib5/url_shortener_PA/config"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/labstack/echo"
)

func BasicAuthDB(username, password string, c echo.Context) (bool, error) {
	var user _model.User
	err := _config.DB.Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
