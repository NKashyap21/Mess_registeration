package services

import (
	"errors"

	"github.com/LambdaIITH/mess_registration/internal/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	A_LDH = iota
	A_UDH
	B_LDH
	B_UDH
)

type SwapService struct {
	db *gorm.DB
}

func NewSwapService(db *gorm.DB) *SwapService {
	return &SwapService{
		db: db,
	}
}

func (s *SwapService) CreateSwapRequest(requesterID uuid.UUID, name, email, requestType string) error {
	var requester schema.User
	if err := s.db.First(&requester, requesterID).Error; err != nil {
		return errors.New("requester not found")
	}

	// check which mess to which mess
	//

	swapRequest := &schema.SwapRequest{
		ID:        uuid.New(),
		Requester: requester,
		Name:      name,
		Email:     email,
		Type:      requestType,
		Status:    "pending",
	}

	return s.db.Create(swapRequest).Error
}

func (s *SwapService) DeleteSwapRequest(requesterID uuid.UUID, name, email string) error {
	result := s.db.Where("requester_id = ? AND  email = ?", requesterID, email).Delete(&schema.SwapRequest{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("swap request not found")
	}
	return nil
}

func (s *SwapService) GetSwapRequestsByUser(userID uuid.UUID) (map[string][]schema.SwapRequest, error) {
	var user schema.User
	if err := s.db.First(&user, userID); err != nil {
		return nil, errors.New("user not found")
	}

	var aToB, bToA []schema.SwapRequest

	// get requests from A to B
	if err := s.db.Preload("Requester").Joins("JOIN users ON users.id = swap_requests.request_id").Where("users.mess = ? AND swap_requests.status = ?", 0, "pending").Find(&aToB).Error; err != nil {
		return nil, err
	}

	if err := s.db.Preload("Requester").Joins("JOIN users ON users.id = swap_requests.request_id").Where("users.mess = ? AND swap_requests.status = ?", 1, "pending").Find(&aToB).Error; err != nil {
		return nil, err
	}

	return map[string][]schema.SwapRequest{
		"A to B": aToB,
		"B to A": bToA,
	}, nil
}

func (s *SwapService) AutoMatchSwaps() error {
	// This is a simplified auto-matching logic
	// In production, you might want more sophisticated matching
	var publicRequestsA []schema.SwapRequest
	var publicRequestsB []schema.SwapRequest

	// Get public requests from mess A
	if err := s.db.Preload("Requester").Joins("JOIN users ON users.id = swap_requests.requester_id").
		Where("users.mess = ? AND swap_requests.type = ? AND swap_requests.status = ?", 0, "public", "pending").
		Find(&publicRequestsA).Error; err != nil {
		return err
	}

	// Get public requests from mess B
	if err := s.db.Preload("Requester").Joins("JOIN users ON users.id = swap_requests.requester_id").
		Where("users.mess = ? AND swap_requests.type = ? AND swap_requests.status = ?", 1, "public", "pending").
		Find(&publicRequestsB).Error; err != nil {
		return err
	}

	// Simple matching - match first available from each side
	minLength := len(publicRequestsA)
	if len(publicRequestsB) < minLength {
		minLength = len(publicRequestsB)
	}

	for i := 0; i < minLength; i++ {
		// Update both requests to approved
		s.db.Model(&publicRequestsA[i]).Update("status", "approved")
		s.db.Model(&publicRequestsB[i]).Update("status", "approved")

		// Swap the mess assignments
		s.db.Model(&publicRequestsA[i].Requester).Update("mess", 1)
		s.db.Model(&publicRequestsB[i].Requester).Update("mess", 0)
	}

	return nil
}
