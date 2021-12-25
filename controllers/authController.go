package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_constant "github.com/fanajib5/url_shortener_PA/constants"
	_m "github.com/fanajib5/url_shortener_PA/middleware"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

// struct to display user data in response api
type UserResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
}

// Check the logged in user whether he is the Super Admin or the Employee
func IsUserLoggedin(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*_m.JwtCustomClaims)
	name := claims.Name
	admin := claims.Admin

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - Success Logged in",
		"data":          name,
		"admin":         admin,
	})
}

// Check if admin is logged in
// or user has permission to customize URL
type UserRole struct {
	Admin     bool `json:"admin"`
	Customize bool `json:"customize"`
}

// validate
func (a *UserRole) ValidateUser(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*_m.JwtCustomClaims)
	isAdmin := claims.Admin
	isCanCustomize := claims.Customize

	// given value that user has logged in as admin
	a.Admin = isAdmin

	// given value that user has logged in as user
	// so they permitted to customize shroten url
	a.Customize = isCanCustomize
}

// function to hash the given password
func HMAC_SHA(rawPwd string) string {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(_constant.SECRET_HMAC))

	// Write Data to it
	h.Write([]byte(rawPwd))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha
}

// --------------------------------------------------------------------------------
//                   REGISTER AND LOGIN SECTION
// --------------------------------------------------------------------------------
// register user
func RegisterUser(c echo.Context) error {
	user := _model.User{}
	c.Bind(&user)

	// validate input data
	validationErrors := validator.New().Struct(user)
	if validationErrors != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       validationErrors.Error(),
		})
	}

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

// login user
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
