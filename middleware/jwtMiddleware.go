package middleware

import (
	"time"

	_constant "github.com/fanajib5/url_shortener_PA/constants"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`
	Customize bool   `json:"customize"`
	jwt.StandardClaims
}

func GenerateToken(userId int, name string, admin bool, customize bool) (string, error) {
	claims := &JwtCustomClaims{
		Id:        userId,
		Name:      name,
		Admin:     admin,
		Customize: customize,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(_constant.SECRET_JWT))
}
