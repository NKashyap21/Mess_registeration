package models

import "time"

type SwapRequest struct {
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Type      string    `json:"type" validate:"required,oneof='friend' 'public'"`
	Password  string    `json:"password" validate:"required,min=6,max=100"`
	UserID    uint      `json:"-" gorm:"foreignKey"`
	Direction string    `json:"direction"` // A to B or B to A
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Swap request only updates when completed
}
