package middleware

import (
	"time"

	_constant "github.com/fanajib5/url_shortener_PA/constants"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func GenerateToken(userId int, name string, admin bool) (string, error) {
	claims := &JwtCustomClaims{
		Id:    userId,
		Name:  name,
		Admin: admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}
	// fmt.Println(claims)
	// claims := jwt.MapClaims{}
	// claims["userId"] = userId
	// claims["name"] = name
	// claims["admin"] = admin
	// claims["exp"] = time.Now().Add(time.Minute * 2).Unix()

	// -----------------
	// GENERATED JWT
	// {
	// 	"claims": {
	// 		"admin": true,
	// 		"exp": 1640272159,
	// 		"id": 1,
	// 		"name": "super admin 01"
	// 	},
	// 	"data": "super admin 01",
	// 	"message": "success logged in"
	// }
	// ----------------------

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(_constant.SECRET_JWT))
}
