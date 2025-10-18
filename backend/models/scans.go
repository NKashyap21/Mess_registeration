package models

import (
	"time"
)

type Scans struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;index"` // foreign key reference to User table
	MessID    uint      `json:"mess_id" gorm:"not null;check:mess_id >= 1 AND mess_id <= 4"`
	Meal      int       `json:"meal" gorm:"not null;check:meal >= 1 AND meal <= 4"`
	Date      time.Time `json:"date" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
