package db

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Phone       string    `json:"phone" gorm:"uniqueIndex" validate:"required,min=10,max=15"`
	RollNo      string    `json:"roll_no" gorm:"uniqueIndex" validate:"required"`
	Mess        int8      `json:"mess" validate:"required,oneof=1 2 3 4 0" default:"0"`
	NextMess    int8      `json:"next_mess" validate:"oneof=1 2 3 4 5 0" default:"0"`
	Type        int8      `json:"type" gorm:"default:0" validate:"oneof=0 1 2"`
	CanRegister bool      `json:"can_register" gorm:"default:true"`
}

// 1 = MessA LDH, 2 = MessA UDH, 3 = MessB LDH, 4 = MessB UDH, 5 = Veg Mess (subset of MessA), 0 = Unassigned
// 0 = Student, 1 = Mess Staff, 2 = Hostel Office
