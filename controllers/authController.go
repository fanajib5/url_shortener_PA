package controller

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type UserResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
}

// Check the logged in user whether he is the Super Admin or the Employee
func IsUserLoggedin(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"]

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - Success Logged in",
		"data":          name,
	})
}

// Check if admin is logged in
// or user has permission to customize URL
type UserRole struct {
	Admin  bool `json:"admin"`
	Custom bool `json:"custom"`
}

// Check if customer is logged in
// type IsCustomer struct {
// 	Custom bool `json:"custom"`
// }

func (a *UserRole) ValidateAdmin(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"]
	isIdValid := claims["id"]

	if isAdmin == true {
		a.Admin = true
	}

	if isIdValid == true {
		a.Custom = true
	}
}
