package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"gorm.io/gorm"
)

type SyncService struct {
	db                    *gorm.DB
	redisService          *RedisMessService
	ticker                *time.Ticker
	stopChan              chan bool
	lastRegistrationState bool
}

func NewSyncService() *SyncService {
	return &SyncService{
		db:           config.GetDB(),
		redisService: NewRedisMessService(),
		stopChan:     make(chan bool),
	}
}

// isRegistrationOpen checks if registration is currently open
func (s *SyncService) isRegistrationOpen() bool {
	var registrationDetails models.MessRegistrationDetails
	if err := s.db.First(&registrationDetails).Error; err != nil {
		return false
	}

	return registrationDetails.NormalRegistrationOpen
}

// StartBackgroundSync starts the background sync process
func (s *SyncService) StartBackgroundSync(intervalSeconds int) {
	s.ticker = time.NewTicker(time.Duration(intervalSeconds) * time.Second)

	go func() {
		log.Printf("Starting background sync service with %d second intervals", intervalSeconds)

		for {
			select {
			case <-s.ticker.C:
				isOpen := s.isRegistrationOpen()

				// Log registration status changes
				if isOpen != s.lastRegistrationState {
					if isOpen {
						log.Println("Registration opened - starting sync operations")
					} else {
						log.Println("Registration closed - stopping sync operations")
					}
					s.lastRegistrationState = isOpen
				}

				if isOpen {
					if err := s.syncPendingRegistrations(); err != nil {
						log.Printf("Error during background sync: %v", err)
					}
				}
			case <-s.stopChan:
				log.Println("Stopping background sync service")
				return
			}
		}
	}()
}

// StopBackgroundSync stops the background sync process
func (s *SyncService) StopBackgroundSync() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.stopChan <- true
}

// syncPendingRegistrations syncs all pending registrations from Redis to database
func (s *SyncService) syncPendingRegistrations() error {
	pendingUsers, err := s.redisService.GetPendingSyncUsers()
	if err != nil {
		return err
	}

	if len(pendingUsers) == 0 {
		return nil // Nothing to sync
	}

	log.Printf("Syncing %d pending registrations to database", len(pendingUsers))

	successCount := 0
	errorCount := 0

	for _, userIDStr := range pendingUsers {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			log.Printf("Invalid user ID format: %s", userIDStr)
			errorCount++
			continue
		}

		if err := s.syncUserRegistration(uint(userID)); err != nil {
			log.Printf("Failed to sync user %d: %v", userID, err)
			errorCount++
		} else {
			successCount++
		}
	}

	log.Printf("Sync completed: %d successful, %d errors", successCount, errorCount)
	return nil
}

// syncUserRegistration syncs a single user's registration from Redis to database
func (s *SyncService) syncUserRegistration(userID uint) error {
	// Get the mess assignment from Redis
	messID, err := s.redisService.GetUserMess(userID)
	if err != nil {
		return err
	}

	if messID == 0 {
		// User has no mess assigned in Redis, remove from pending sync
		return s.redisService.RemoveFromPendingSync(userID)
	}

	// Begin database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user from database
	var user db.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			log.Printf("User %d not found in database, removing from Redis", userID)
			// User doesn't exist in DB, clean up Redis
			s.redisService.ClearUserRegistration(userID, messID)
			return nil
		}
		return err
	}

	// Check if user already has a different mess assigned in NextMess
	if user.NextMess != 0 && user.NextMess != int8(messID) {
		tx.Rollback()
		log.Printf("Conflict: User %d has NextMess %d in DB but %d in Redis", userID, user.NextMess, messID)
		// Remove from Redis to avoid further conflicts
		s.redisService.ClearUserRegistration(userID, messID)
		return nil
	}

	// Update user's NextMess in database
	user.NextMess = int8(messID)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Remove from pending sync queue
	if err := s.redisService.RemoveFromPendingSync(userID); err != nil {
		log.Printf("ERROR: Failed to remove user %d from pending sync queue: %v", userID, err)
		// Return error to retry sync later
		return fmt.Errorf("failed to remove from pending sync: %w", err)
	}

	log.Printf("Successfully synced user %d with mess %d to database", userID, messID)
	return nil
}

// InitializeRedisFromDB initializes Redis counters and capacities from database state
func (s *SyncService) InitializeRedisFromDB() error {
	log.Println("Initializing Redis from database state...")

	// Get registration details to set capacities
	var regDetails models.MessRegistrationDetails
	if err := s.db.First(&regDetails).Error; err != nil {
		return err
	}

	// Set capacities in Redis
	capacities := map[int]int{
		1: regDetails.MessALDHCapacity,
		2: regDetails.MessAUDHCapacity,
		3: regDetails.MessBLDHCapacity,
		4: regDetails.MessBUDHCapacity,
		5: regDetails.MessALDHCapacity + regDetails.MessAUDHCapacity, // Veg mess shares capacity
	}

	if err := s.redisService.InitializeMessCapacities(capacities); err != nil {
		return err
	}

	// Count current registrations in database and sync to Redis
	// Using NextMess column as that's where new registrations are stored
	for messID := 1; messID <= 5; messID++ {
		var count int64
		if err := s.db.Model(&models.User{}).Where("next_mess = ?", messID).Count(&count).Error; err != nil {
			return err
		}

		// Set the count in Redis (this will override any existing count)
		counterKey := fmt.Sprintf("mess:%d:count", messID)
		if err := s.redisService.client.Set(s.redisService.ctx, counterKey, count, 0).Err(); err != nil {
			return err
		}

		log.Printf("Initialized mess %d: %d/%d registrations", messID, count, capacities[messID])
	}

	// Get all users with NextMess assignments and sync to Redis
	var users []models.User
	if err := s.db.Where("next_mess > 0").Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		userMessKey := fmt.Sprintf("user:%d:mess", user.ID)
		if err := s.redisService.client.Set(s.redisService.ctx, userMessKey, user.NextMess, 0).Err(); err != nil {
			log.Printf("Warning: Failed to sync user %d to Redis: %v", user.ID, err)
		}
	}

	log.Printf("Redis initialization complete. Synced %d user assignments", len(users))
	return nil
}
