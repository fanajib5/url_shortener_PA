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
	eJwt.Use(echoMiddleware.JWT([]byte(_constant.SECRET_JWT)))
	eJwt.GET("/", _controller.IsUserLoggedin)

	// eJwt.GET("/orders", controller.GetOrders)
	// eJwt.POST("/orders", controller.CreateOrder)
	// eJwt.PUT("/orders/:id", controller.UpdateOrder)
	// eJwt.DELETE("/orders/:id", controller.DeleteOrder)

	// eJwt.GET("/cars", controller.GetCars)
	// eJwt.POST("/cars", controller.CreateCar)
	// eJwt.PUT("/cars/:id", controller.UpdateCar)
	// eJwt.DELETE("/cars/:id", controller.DeleteCar)

	return e
}
