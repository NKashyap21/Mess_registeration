package services

import (
	"errors"

	"github.com/LambdaIITH/mess_registration/internal/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(user *schema.User) error {
	var existingUser schema.User

	if err := s.db.Where("user_id = ? OR email = ? OR roll_no = ?", user.UserID, user.Email, user.RollNo).First(&existingUser).Error; err != nil {
		return errors.New("user already exists with this user_id, email or roll_no")
	}

	if user.Mess < 0 || user.Mess > 3 {
		return errors.New("invalid mess number")
	}

	if user.UserType == "" {
		user.UserType = "student"
	}

	return s.db.Create(user).Error
}

func (s *UserService) GetUserByUserID(userID string) (*schema.User, error) {
	var user schema.User

	if err := s.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*schema.User, error) {
	var user schema.User

	if err := s.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByRollNo(rollNo string) (*schema.User, error) {
	var user schema.User
	err := s.db.Where("roll_no = ?", rollNo).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&schema.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) GetAllUsers() ([]schema.User, error) {
	var users []schema.User
	err := s.db.Find(&users).Error
	return users, err
}

func (s *UserService) DeleteUser(userId uuid.UUID) error {
	return s.db.Delete(&schema.User{}, "user_id = ?", userId).Error
}

func (s *UserService) ResetAllUsers() error {
	return s.db.Exec("DELETE FROM users").Error
}
