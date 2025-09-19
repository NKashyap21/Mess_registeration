package models

type SwapRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Type      string `json:"type" validate:"required,oneof='friend' 'public'"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
	UserID    uint   `json:"-" gorm:"foreignKey"`
	Direction string `json:"-"` // A to B or B to A
}
