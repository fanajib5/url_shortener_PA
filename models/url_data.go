package model

import "gorm.io/gorm"

// UUID       string `json:"uuid" form:"uuid"`

type UrlData struct {
	gorm.Model
	Custom     bool   `json:"custom" form:"custom"`
	Url        string `json:"url" form:"url"`
	ShortUrl   string `json:"short_url" form:"short_url"`
	HitCounter int    `json:"hit_counter" form:"hit_counter"`
}
