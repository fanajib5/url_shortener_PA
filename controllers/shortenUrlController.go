package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_config "github.com/fanajib5/url_shortener_PA/config"
	_m "github.com/fanajib5/url_shortener_PA/middleware"
	_model "github.com/fanajib5/url_shortener_PA/models"

	"github.com/labstack/echo"
	nanoid "github.com/matoous/go-nanoid/v2"
)

//
// we intended to simplify API responses by using these two structs
// feel free to correct me if you guys have the better approach :D
//
type UrlCollection struct {
	ActualUrl   string `json:"actual_url"`
	ShortUrl    string `json:"short_url"`
	RedirectUrl string `json:"redirect_url"`
}

type UserData struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`
	Customize bool   `json:"customize"`
}

func (udt *UserData) ParseJWT(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*_m.JwtCustomClaims)
	udt.Id = claims.Id
	udt.Name = claims.Name
	udt.Admin = claims.Admin
	udt.Customize = claims.Customize
}

// create short url for guest (not-logged user)
func CreateShortUrl(c echo.Context) error {
	urlData := _model.UrlData{}
	c.Bind(&urlData)
	u := &UrlCollection{}
	env := &_config.DotEnv{}
	env.LoadDotEnv()

	nanoId, errNanoId := nanoid.Generate(env.SECRET_NANOID, 6)
	if errNanoId != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response_code": int(http.StatusInternalServerError),
			"message":       errNanoId.Error(),
		})
	}

	if urlData.ShortUrl != "" {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"response_code": int(http.StatusConflict),
			"message":       "if you want to customize short url, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	urlData.ShortUrl = c.Request().Host + "/" + nanoId
	urlData.CreatedBy = 0
	urlData.UpdatedBy = 0

	errStoreUrl := _config.DB.Save(&urlData).Error
	if errStoreUrl != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response_code": int(http.StatusInternalServerError),
			"message":       errStoreUrl.Error(),
		})
	}

	u.ActualUrl = urlData.Url
	u.ShortUrl = urlData.ShortUrl

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - URL clipped",
		"data":          u,
	})
}

// END OF create short url for guest (not-logged user)

// redirect short url to the actual url (all user whether logged user or not)
func RedirectClippedUrl(c echo.Context) error {
	urlData := _model.UrlData{}
	c.Bind(&urlData)
	u := &UrlCollection{}

	shortUrl := c.Param("urlShortCode")

	errSearchedUrl := _config.DB.Model(&urlData).Where("short_url like ? and deleted_at is NULL", "%"+shortUrl).Find(&urlData).Error
	if errSearchedUrl != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchedUrl.Error(),
		})
	}

	// we add counter by 1 whenever the short_url fired
	counterAddition := urlData.HitCounter + 1
	updateHitCounter := _config.DB.Model(&urlData).Where("short_url like ? and deleted_at is NULL", "%"+shortUrl).UpdateColumns(map[string]interface{}{
		"hit_counter": counterAddition,
	})

	if updateHitCounter.Error != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       updateHitCounter.Error.Error(),
		})
	}

	// print hit_counter just to make sure that counterAddition works well hehehe
	// fmt.Println("hit_counter:", urlData.HitCounter)

	u.ActualUrl = urlData.Url
	u.ShortUrl = urlData.ShortUrl
	u.RedirectUrl = urlData.Url

	//
	// user can directly go to the origin url by hit redirect_url
	// from the frontend route management CMIIW
	//

	return c.JSON(http.StatusFound, map[string]interface{}{
		"response_code": int(http.StatusFound),
		"message":       "OK - Clipped URL found",
		"data":          u,
	})
}

// END OF redirect short url to the actual url (all user whether logged user or not)

// ------------------------------------------------------------------------------------------------
//                            this is logged user section, thanks :D
// ------------------------------------------------------------------------------------------------
// create custom url for logged user
func CustomShortUrl(c echo.Context) error {
	a := &UserRole{}
	a.ValidateAdmin(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to customize clipped URL, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	// token := &jwt.Token{}
	// if !token.Valid {
	// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
	// 		"response_code": int(http.StatusUnauthorized),
	// 		"message":       "It seems your token has been invalid, please login again!",
	// 		"register_url":  c.Request().Host + "/register",
	// 		"login_url":     c.Request().Host + "/login",
	// 	})
	// }

	urlData := _model.UrlData{}
	c.Bind(&urlData)
	u := &UrlCollection{}

	udt := &UserData{}
	udt.ParseJWT(c)

	buffUrl := urlData.ShortUrl
	urlData.ShortUrl = c.Request().Host + "/" + buffUrl
	urlData.Custom = true
	urlData.CreatedBy = udt.Id
	urlData.UpdatedBy = udt.Id

	errStoreUrl := _config.DB.Model(&urlData).FirstOrInit(&urlData).Error
	if errStoreUrl != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response_code": int(http.StatusInternalServerError),
			"message":       errStoreUrl.Error(),
		})
	}

	u.ActualUrl = urlData.Url
	u.ShortUrl = urlData.ShortUrl

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - URL clipped and customized",
		"data":          u,
	})
}

// END OF create custom url for logged user

// get all clipped url which created by user
func GetShortUrlList(c echo.Context) error {
	a := &UserRole{}
	a.ValidateAdmin(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to get clipped URL list, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	urlData := _model.UrlData{}
	c.Bind(&urlData)

	udt := UserData{}
	udt.ParseJWT(c)

	errSearchedUrl := _config.DB.Model(&urlData).Where("created_by = ? and deleted_at is NULL", udt.Id).Find(&urlData).Error
	if errSearchedUrl != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchedUrl.Error() + " - hmmm... it seems that you are didn't using our service yet! :D",
		})
	}

	var count int64
	errCountUrlList := _config.DB.Model(&urlData).Where("created_by = ? and deleted_at is NULL", udt.Id).Group("created_by").Count(&count).Error
	if errCountUrlList != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errCountUrlList.Error() + " - hmmm... it seems that you are didn't using our service yet! :D",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK",
		"data":          urlData,
		"records_total": count,
	})
}

// END OF get all clipped url which created by user

// get all clipped url which created by user
func GetShortUrl(c echo.Context) error {
	a := &UserRole{}
	a.ValidateAdmin(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to get clipped URL data, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	urlData := _model.UrlData{}
	urlId, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&urlData)

	udt := UserData{}
	udt.ParseJWT(c)

	errSearchedUrl := _config.DB.Model(&urlData).Where("id = ?  and deleted_at is NULL", urlId).Find(&urlData).Error
	if errSearchedUrl != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchedUrl.Error() + " - hmmm... it seems that you are didn't using our service yet! :D",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - clipped url data found",
		"data":          urlData,
	})
}

// END OF get all clipped url which created by user

// get all clipped url which created by user
func UpdateShortUrl(c echo.Context) error {
	a := &UserRole{}
	a.ValidateAdmin(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to update the clipped URL data, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	urlData := _model.UrlData{}
	urlId, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&urlData)

	udt := UserData{}
	udt.ParseJWT(c)

	// we checking url data before user make changes to their data :D
	checkChangedUrlData := _config.DB.Model(&urlData).Where("short_url = ?", urlData.ShortUrl).Find(&urlData).Error
	if checkChangedUrlData == nil {
		return c.JSON(http.StatusFound, map[string]interface{}{
			"response_code": int(http.StatusFound),
			"message":       checkChangedUrlData.Error() + " - Ooppss! Your new short url has been used by other user. PLease use another one :D",
		})
	}

	updateUrlData := _config.DB.Model(&urlData).Where("id = ?  and deleted_at is NULL", urlId).Find(&urlData).UpdateColumns(map[string]interface{}{
		"short_url": urlData.ShortUrl,
	})

	if updateUrlData.Error != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"response_code": int(http.StatusConflict),
			"message":       updateUrlData.Error.Error() + " - it seems that your input data conflict with our specification :D",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response_code": int(http.StatusOK),
		"message":       "OK - url data updated",
		"data":          urlData,
	})
}

// END OF get all clipped url which created by user
func DeleteShortUrl(c echo.Context) error {
	a := &UserRole{}
	a.ValidateAdmin(c)
	if !a.Customize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"response_code": int(http.StatusUnauthorized),
			"message":       "You don't have permission to update the clipped URL data, please login or register",
			"register_url":  c.Request().Host + "/register",
			"login_url":     c.Request().Host + "/login",
		})
	}

	urlData := _model.UrlData{}
	urlId, _ := strconv.Atoi(c.Param("id"))
	c.Bind(&urlData)

	udt := UserData{}
	udt.ParseJWT(c)

	// soft delete procedure
	deleteData := _config.DB.Model(&urlData).Where("id = ?  and deleted_at is NULL", urlId).Find(&urlData).UpdateColumns(map[string]interface{}{
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
