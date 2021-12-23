package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_m "github.com/fanajib5/url_shortener_PA/middleware"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/labstack/echo"
)

func RegisterUser(c echo.Context) error {
	user := _model.User{}
	c.Bind(&user)

	rawPass := user.Password
	h := sha256.New()
	h.Write([]byte(rawPass))
	user.Password = hex.EncodeToString(h.Sum(nil))

	errRegister := _config.DB.Save(&user).Error
	if errRegister != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"message":     "failed login",
			"description": errRegister.Error(),
		})
	}

	token, errGen := _m.GenerateToken(user.Id, user.Fullname, user.Admin)
	if errGen != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":     "failed login",
			"description": errGen.Error(),
		})
	}

	generatedToken := &jwt.Token{}
	claims := generatedToken.Claims.(jwt.MapClaims)
	expAt := claims["exp"]

	userResponse := UserResponse{
		Username: user.Username,
		Name:     user.Fullname,
		Admin:    user.Admin,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success login",
		"user_data":  userResponse,
		"token":      token,
		"expired_at": expAt,
	})
}

func LoginUser(c echo.Context) error {
	user := _model.User{}
	c.Bind(&user)

	errLogin := _config.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error
	if errLogin != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"message":     "failed login",
			"description": errLogin.Error(),
		})
	}

	token, errGen := _m.GenerateToken(user.Id, user.Fullname, user.Admin)
	if errGen != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":     "failed login",
			"description": errGen.Error(),
		})
	}

	generatedToken := &jwt.Token{}
	claims := generatedToken.Claims.(jwt.MapClaims)
	expAt := claims["exp"]

	userResponse := UserResponse{
		Username: user.Username,
		Name:     user.Fullname,
		Admin:    user.Admin,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success login",
		"user_data":  userResponse,
		"token":      token,
		"expired_at": expAt,
	})
}
