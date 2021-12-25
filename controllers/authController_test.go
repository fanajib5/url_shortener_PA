package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_controller "github.com/fanajib5/url_shortener_PA/controllers"
	model "github.com/fanajib5/url_shortener_PA/models"
	_routes "github.com/fanajib5/url_shortener_PA/routes"

	"github.com/gavv/httpexpect"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func PrepareTestData() error {
	_config.InitDB()

	buffer_pwd := _controller.HMAC_SHA("unit-test")

	_config.DB.Create(model.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Id:       99,
		Fullname: "unit-test",
		Email:    "unit-test",
		Username: "unit-test",
		Password: buffer_pwd,
		Admin:    true,
	})

	return nil
}

func InitEcho() *echo.Echo {

	e := _routes.New()

	return e
}

func TestIsUserLoggedin(t *testing.T) {

	PrepareTestData()

	handler := InitEcho()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect for Mock instance
	e := httpexpect.New(t, server.URL)

	//==================GET JWT TOKEN FOR ADD IN HEADER REQUEST===================
	data := map[string]interface{}{
		"username": "unit-test",
		"password": "unit-test",
	}

	// get token
	obj := e.POST("/login").WithJSON(data).Expect().
		Status(http.StatusOK).JSON().Object()

	token := obj.Value("token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	t.Run("Expected Get validation the loggedin user", func(t *testing.T) {
		auth.GET("/cust/").
			Expect().
			Status(http.StatusOK).JSON().Object()
	})
}

func TestLoginUser(t *testing.T) {

	// PrepareTestData()

	handler := InitEcho()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect for Mock instance
	e := httpexpect.New(t, server.URL)

	// buffer_pwd := _controller.HMAC_SHA("unit-test2")

	t.Run("Expected Login user", func(t *testing.T) {
		data := map[string]interface{}{
			"username": "unit-test",
			"password": "unit-test",
		}

		e.POST("/login").WithJSON(data).
			Expect().Status(http.StatusOK)
	})

}

func TestRegisterUser(t *testing.T) {

	// PrepareTestData()

	handler := InitEcho()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect for Mock instance
	e := httpexpect.New(t, server.URL)

	// buffer_pwd := _controller.HMAC_SHA("unit-test2")

	t.Run("Expected Insert user", func(t *testing.T) {
		data := map[string]interface{}{
			"username": "unit-test",
			"password": "unit-test",
		}

		e.POST("/register").WithJSON(data).
			Expect().Status(http.StatusUnprocessableEntity)
	})

}
