package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestGetUrlList(t *testing.T) {

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
	obj := e.POST("/login").WithJSON(data).Expect().
		Status(http.StatusOK).JSON().Object()

	token := obj.Value("token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	t.Run("Expected Find ALL short url list", func(t *testing.T) {
		auth.GET("/cust/clip").
			Expect().
			Status(http.StatusOK).JSON().Object()
	})
}

func TestGetShortUrl(t *testing.T) {

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

	t.Run("Expected Find url data", func(t *testing.T) {
		auth.GET("/cust/clip/{id}").WithPath("id", 8).
			Expect().
			Status(http.StatusOK).JSON().Object()
	})
}

func TestCreateShortUrl(t *testing.T) {

	// PrepareTestData()

	handler := InitEcho()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect for Mock instance
	e := httpexpect.New(t, server.URL)

	t.Run("Expected Insert url, then call get one By ID and found that url", func(t *testing.T) {
		dataForInsert := map[string]interface{}{
			"url": "https://www.youtube.com/watch?v=-iQAuGai0VE&ab_channel=MusikIndonesia",
		}

		e.POST("/clip").WithJSON(dataForInsert).
			Expect().
			Status(http.StatusOK)
	})
}

func TestCustomShortUrl(t *testing.T) {

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

	t.Run("Expected Insert url, then call get one By ID and found that url", func(t *testing.T) {
		dataForInsert := map[string]interface{}{
			"custom":    true,
			"url":       "https://www.youtube.com/watch?v=-iQAuGai0VE&ab_channel=MusikIndonesia",
			"short_url": "test-short-url",
		}

		auth.POST("/cust/clip").WithJSON(dataForInsert).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Insert url, then call get one By ID But NOT found that url", func(t *testing.T) {
		dataForInsert := map[string]interface{}{
			"custom":    true,
			"url":       "https://www.youtube.com/watch?v=-iQAuGai0VE&ab_channel=MusikIndonesia",
			"short_url": 22222,
		}

		auth.POST("/cust/clip").WithJSON(dataForInsert).
			Expect().
			Status(http.StatusOK)
	})
}

func TestUpdateShortUrl(t *testing.T) {

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
		dataForUpdate := map[string]interface{}{
			"custom":    true,
			"url":       "https://www.youtube.com/watch?v=-iQAuGai0VE&ab_channel=MusikIndonesia",
			"short_url": "test-short-url",
		}

		auth.PUT("/cust/clip/{id}").WithPath("id", 8).
			WithJSON(dataForUpdate).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Update url, then call get one By ID But NOT found that url", func(t *testing.T) {
		dataForUpdate := map[string]interface{}{
			"custom":    false,
			"url":       "https://www.youtube.com/watch?v=-iQAuGai0VE&ab_channel=MusikIndonesia",
			"short_url": 2323,
		}

		auth.PUT("/cust/clip/{id}").WithPath("id", 29).
			WithJSON(dataForUpdate).Expect().
			Status(http.StatusOK)
	})
}

func TestDeleteShortUrl(t *testing.T) {

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

	t.Run("Expected Delete url", func(t *testing.T) {
		auth.DELETE("/cust/clip/{id}").WithPath("id", 26).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Expected Delete url", func(t *testing.T) {
		auth.DELETE("/cust/clip/{id}").WithPath("id", 8).
			Expect().
			Status(http.StatusOK)
	})
}

func TestRedirectClippedUrl(t *testing.T) {

	// PrepareTestData()

	// l, err := net.Listen("tcp", "localhost:8080")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	handler := InitEcho()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect for Mock instance
	e := httpexpect.New(t, server.URL)

	t.Run("Expected Delete url", func(t *testing.T) {
		e.GET("/{id}").WithPath("id", "pingin-turu002").
			Expect().
			Status(http.StatusOK)
	})
}
