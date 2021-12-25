package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username" gorm:"unique" validate:"required,alphanumunicode"`
	Fullname string `json:"fullname" form: "fullname" validate:"required,alphanumunicode"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,gte=8,lte=32"`
	Admin    bool   `json:"admin" form: "admin" gorm:"default:false" validate:"boolean"`
}
