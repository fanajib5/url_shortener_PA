package model

import "gorm.io/gorm"

type UrlData struct {
	gorm.Model
	Custom     bool   `json:"custom" form:"custom" gorm:"default:false"`
	Url        string `json:"url" form:"url" validate:"required,url"`
	ShortUrl   string `json:"short_url" form:"short_url" gorm:"unique"`
	HitCounter int    `json:"hit_counter" form:"hit_counter"`
	CreatedBy  int    `json:"created_by" form:"created_by"`
	UpdatedBy  int    `json:"updated_by" form:"updated_by"`
	DeletedBy  int    `json:"deleted_by" form:"deleted_by"`
}
