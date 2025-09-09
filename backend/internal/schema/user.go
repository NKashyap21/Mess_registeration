package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"unique;not null" binding:"required"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	RollNo    string    `json:"rollno" gorm:"unique;not null"`
	UserType  string    `json:"user_type" gorm:"not null"` // student, admin, mess_staff
	VegType   string    `json:"veg_type" gorm:"not null"`  // veg, non-veg
	Mess      int       `json:"mess" gorm:"not null"`      // 0, 1, 2, 3
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SwapRequest struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RequesterID uuid.UUID `json:"requester_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null" binding:"required"`
	Email       string    `json:"email" gorm:"not null" binding:"required,email"`
	Type        string    `json:"type" gorm:"not null" binding:"required,oneof=friend public"` // friend, public
	Status      string    `json:"status" gorm:"default:pending"`                               // pending, approved, rejected
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign key relationship
	Requester User `json:"requester" gorm:"foreignKey:RequesterID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (sr *SwapRequest) BeforeCreate(tx *gorm.DB) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	return nil
}
