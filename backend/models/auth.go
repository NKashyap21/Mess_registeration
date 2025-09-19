package models

type GoogleLoginData struct {
	Token string `json:"token" binding:"required"`
}