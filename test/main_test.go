package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	controller "github.com/fanajib5/url_shortener_PA/controllers"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var userJSON = `{"username":"faiq123","password":"faiq123"}`

func TestLogin(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(userJSON))
	fmt.Println(userJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// h := &handler{mockDB}

	// Assertions
	if assert.NoError(t, controller.LoginUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}
}
