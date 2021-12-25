package routes

import (
	_constant "github.com/fanajib5/url_shortener_PA/constants"
	_controller "github.com/fanajib5/url_shortener_PA/controllers"
	_m "github.com/fanajib5/url_shortener_PA/middleware"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	_m.LogMiddleware(e)
	e.POST("/clip", _controller.CreateShortUrl)
	e.GET("/:urlShortCode", _controller.RedirectClippedUrl)

	e.POST("/register", _controller.RegisterUser)
	e.POST("/login", _controller.LoginUser)

	eJwt := e.Group("/cust")
	config := echoMiddleware.JWTConfig{
		Claims:     &_m.JwtCustomClaims{},
		SigningKey: []byte(_constant.SECRET_JWT),
	}
	eJwt.Use(echoMiddleware.JWTWithConfig(config))
	// eJwt.Use(echoMiddleware.JWT([]byte(_constant.SECRET_JWT)))
	eJwt.GET("/", _controller.IsUserLoggedin)

	eJwt.POST("/clip", _controller.CustomShortUrl)
	eJwt.GET("/clip", _controller.GetShortUrlList)
	eJwt.GET("/clip/:id", _controller.GetShortUrl)
	eJwt.PUT("/clip/:id", _controller.UpdateShortUrl)
	eJwt.DELETE("/clip/:id", _controller.DeleteShortUrl)

	eJwt.GET("/profile/:id", _controller.GetUserData)
	eJwt.PUT("/profile/:id", _controller.UpdateUserData)
	eJwt.DELETE("/profile/:id", _controller.DeleteUser)

	return e
}
