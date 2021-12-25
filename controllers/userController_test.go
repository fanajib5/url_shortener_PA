package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_controller "github.com/fanajib5/url_shortener_PA/controllers"

	"github.com/gavv/httpexpect"
)

func TestGetUserData(t *testing.T) {

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
	obj := e.POST("/login").
		WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token := obj.Value("token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	t.Run("Expected Find user data", func(t *testing.T) {
		auth.GET("/cust/profile/{id}").WithPath("id", 100).
			Expect().
			Status(http.StatusOK).JSON().Object()
	})

	t.Run("Expected Find user data", func(t *testing.T) {
		auth.GET("/profile/{id}").WithPath("id", 5).
			Expect().
			Status(http.StatusNotFound).JSON().Object()
	})
}

func TestUpdateUserData(t *testing.T) {

	// PrepareTestData()

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
	obj := e.POST("/login").
		WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token := obj.Value("token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	t.Run("Expected Update url, then call get one By ID and found that url", func(t *testing.T) {
		buffer_pwd := _controller.HMAC_SHA("unit-test")

		dataForUpdate := map[string]interface{}{
			"fullname": "unit test tes",
			"email":    "unit-test",
			"username": "unit-test",
			"password": buffer_pwd,
			"admin":    true,
		}

		auth.PUT("/cust/profile/{id}").WithPath("id", 101).
			WithJSON(dataForUpdate).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Update url, then call get one By ID But NOT found that url", func(t *testing.T) {
		buffer_pwd := _controller.HMAC_SHA("unit-test")

		dataForUpdate := map[string]interface{}{
			"fullname": "unit test tes",
			"email":    "unit-test",
			"username": "unit-test",
			"password": buffer_pwd,
			"admin":    true,
		}

		auth.PUT("/cust/profile/{id}").WithPath("id", 29).
			WithJSON(dataForUpdate).Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Update url, then call get one By ID But NOT found that url", func(t *testing.T) {
		buffer_pwd := _controller.HMAC_SHA("unit-test")

		dataForUpdate := map[string]interface{}{
			"fullname": "unit test tes",
			"email":    "unit-test",
			"username": "unit-test",
			"password": buffer_pwd,
			"admin":    true,
		}

		auth.PUT("/profile/{id}").WithPath("id", 29).
			WithJSON(dataForUpdate).Expect().
			Status(http.StatusNotFound)
	})
}

func TestDeleteUserData(t *testing.T) {

	// PrepareTestData()

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
	obj := e.POST("/login").
		WithJSON(data).
		Expect().
		Status(http.StatusOK).JSON().Object()

	token := obj.Value("token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	t.Run("Expected Update url, then call get one By ID and found that url", func(t *testing.T) {
		auth.DELETE("/cust/profile/{id}").WithPath("id", 101).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Update url, then call get one By ID But NOT found that url", func(t *testing.T) {
		auth.DELETE("/profile/{id}").WithPath("id", 29).
			Expect().
			Status(http.StatusNotFound)
	})
}
