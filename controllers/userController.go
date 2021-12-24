package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_m "github.com/fanajib5/url_shortener_PA/middleware"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/labstack/echo"
)

func HMAC_SHA(rawPwd string) string {
	env := &_config.DotEnv{}
	env.LoadDotEnv()

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(env.SECRET_HMAC))

	// Write Data to it
	h.Write([]byte(rawPwd))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha
}

func RegisterUser(c echo.Context) error {
	user := _model.User{}
	c.Bind(&user)

	// store user password in temporary variable, then hashing it
	buffer_pwd := user.Password
	user.Password = HMAC_SHA(buffer_pwd)

	errRegister := _config.DB.Save(&user).Error
	if errRegister != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       errRegister.Error(),
		})
	}

	userResponse := &UserResponse{
		Username: user.Username,
		Name:     user.Fullname,
		Admin:    user.Admin,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - Success Register",
		"user_data":     userResponse,
	})
}

func LoginUser(c echo.Context) error {
	user := _model.User{}
	c.Bind(&user)

	// hashing user password
	shaPwd := HMAC_SHA(user.Password)

	errLogin := _config.DB.Where("username = ? AND password = ?", user.Username, shaPwd).First(&user).Error
	if errLogin != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       errLogin.Error(),
		})
	}

	storedPwd := user.Password
	if !hmac.Equal([]byte(shaPwd), []byte(storedPwd)) {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       errLogin.Error(),
		})
	}

	token, errGen := _m.GenerateToken(user.Id, user.Fullname, user.Admin, true)
	if errGen != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response_code": int(http.StatusInternalServerError),
			"message":       errLogin.Error(),
		})
	}

	// x := &_m.JwtCustomClaims{}

	userResponse := &UserResponse{
		Username: user.Username,
		Name:     user.Fullname,
		Admin:    user.Admin,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code":    int(http.StatusOK),
		"message":          "OK - Success Login",
		"user_data":        userResponse,
		"token":            token,
		"token_expires_at": 300,
	})
}
