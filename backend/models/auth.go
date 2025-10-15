package models

type GoogleLoginData struct {
	Code string `form:"code" binding:"required"`
}

