package controller

import (
	"fmt"
	"net/http"

	_config "github.com/fanajib5/url_shortener_PA/config"
	_constant "github.com/fanajib5/url_shortener_PA/constants"
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

func CreateShortUrl(c echo.Context) error {
	urlData := _model.UrlData{}
	c.Bind(&urlData)
	u := UrlCollection{}

	nanoId, errNanoId := nanoid.Generate(_constant.CustomAlphanumeric, 6)
	if errNanoId != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"response_code": int(http.StatusInternalServerError),
			"message":       errNanoId.Error(),
		})
	}

	urlData.ShortUrl = c.Request().Host + "/" + nanoId

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

func RedirectClippedUrl(c echo.Context) error {
	urlData := _model.UrlData{}
	c.Bind(&urlData)
	u := UrlCollection{}

	shortUrl := c.Param("urlShortCode")

	errSearchedUrl := _config.DB.Model(&urlData).Where("short_url like ?", "%"+shortUrl).Find(&urlData).Error
	if errSearchedUrl != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"response_code": int(http.StatusNotFound),
			"message":       errSearchedUrl.Error(),
		})
	}

	// we add counter by 1 whenever the short_url fired
	counterAddition := urlData.HitCounter + 1
	updateHitCounter := _config.DB.Model(&urlData).Where("short_url like ?", "%"+shortUrl).UpdateColumns(map[string]interface{}{
		"hit_counter": counterAddition,
	})

	if updateHitCounter.Error != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"response_code": int(http.StatusUnprocessableEntity),
			"message":       updateHitCounter.Error.Error(),
		})
	}

	// print hit_counter just to make sure that counterAddition works well hehehe
	fmt.Println("hit_counter:", urlData.HitCounter)

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
