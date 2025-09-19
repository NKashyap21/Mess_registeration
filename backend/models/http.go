package models

type APIResponse struct {
	Message string      `json:"message" omitempty:"true"`
	Data    interface{} `json:"data" omitempty:"true"`
}
