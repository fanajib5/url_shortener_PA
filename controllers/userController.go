package controller

import (
	"net/http"
	"strconv"
	"time"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/labstack/echo"
)

// get all clipped url which created by user
func GetUserData(c echo.Context) error {
	a := &UserRole{}
	a.ValidateUser(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to get clipped URL data, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	user := _model.User{}
	userId, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&user)

	errSearchUserData := _config.DB.Model(&user).Where("id = ? and deleted_at is NULL", userId).Find(&user).Error
	if errSearchUserData != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchUserData.Error() + " - hmmm... it seems that your profile deleted or you did'not register yet! :D",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - user data found",
		"data":          user,
	})
}

// get all clipped url which created by user
func UpdateUserData(c echo.Context) error {
	a := &UserRole{}
	a.ValidateUser(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to update the clipped URL data, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	user := _model.User{}
	userId, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&user)

	updateUrlData := _config.DB.Model(&user).Where("id = ? and deleted_at is NULL", userId).UpdateColumns(map[string]interface{}{
		"fullname": user.Fullname,
		"email":    user.Email,
	})

	if updateUrlData.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"response_code": int(http.StatusConflict),
			"message":       updateUrlData.Error.Error() + " - it seems that your input data conflict with our specification :D",
		})
	}

	errSearchedUrl := _config.DB.Model(&user).Where("id = ? and deleted_at is NULL", userId).Find(&user).Error
	if errSearchedUrl != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchedUrl.Error() + " - hmmm... it seems that you are didn't using our service yet! :D",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - url data updated",
		"data":          user,
	})
}

// END OF get all clipped url which created by user
func DeleteUser(c echo.Context) error {
	a := &UserRole{}
	a.ValidateUser(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to delete user, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	user := _model.User{}
	userID, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&user)

	// soft delete procedure
	deleteData := _config.DB.Model(&user).Where("id = ? and deleted_at is NULL", userID).Find(&user).UpdateColumns(map[string]interface{}{
		"deleted_at": time.Now(),
	})

	if deleteData.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       deleteData.Error.Error() + " - it seems that your data has deleted :D",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - url data deleted succesfully",
	})
}
