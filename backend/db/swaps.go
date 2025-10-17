package db

import "time"

type SwapRequest struct {
	Type      string    `json:"type" validate:"required,oneof='friend' 'public'"`
	Password  string    `json:"password" validate:"required,min=6,max=100"`
	UserID    uint      `json:"-" gorm:"foreign,primaryKey"`
	Direction string    `json:"-"` // A to B or B to A
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Swap request only updates when completed
}
