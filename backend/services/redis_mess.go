package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/redis/go-redis/v9"
)

type RedisMessService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisMessService() *RedisMessService {
	return &RedisMessService{
		client: config.GetRedisClient(),
		ctx:    context.Background(),
	}
}

// GetMessName returns a human-readable name for the mess ID
func GetMessName(messID int) string {
	switch messID {
	case MESS_A_LDH:
		return "MessA LDH"
	case MESS_A_UDH:
		return "MessA UDH"
	case MESS_B_LDH:
		return "MessB LDH"
	case MESS_B_UDH:
		return "MessB UDH"
	case VEG_MESS:
		return "Veg Mess"
	case NO_MESS:
		return "No Mess"
	default:
		return fmt.Sprintf("Unknown Mess (%d)", messID)
	}
}

// IsValidMessID checks if the mess ID is valid (1-5)
func IsValidMessID(messID int) bool {
	return messID >= MESS_A_LDH && messID <= VEG_MESS
}

// Redis key patterns
const (
	MESS_COUNTER_KEY  = "mess:%d:count"    // mess:1:count, mess:2:count, etc.
	USER_MESS_KEY     = "user:%d:mess"     // user:123:mess -> mess_id
	PENDING_SYNC_KEY  = "pending_sync"     // Set of user IDs that need DB sync
	MESS_CAPACITY_KEY = "mess:%d:capacity" // mess:1:capacity -> max capacity
)

// Mess mapping constants
const (
	MESS_A_LDH = 1 // MessA Lower Dining Hall
	MESS_A_UDH = 2 // MessA Upper Dining Hall
	MESS_B_LDH = 3 // MessB Lower Dining Hall
	MESS_B_UDH = 4 // MessB Upper Dining Hall
	VEG_MESS   = 5 // Vegetarian Mess
	NO_MESS    = 0 // No mess assigned
)

// LoadCapacitiesFromDB loads capacity limits from database and returns them as a map
func (r *RedisMessService) LoadCapacitiesFromDB() (map[int]int, error) {
	db := config.GetDB()

	var regDetails models.MessRegistrationDetails
	if err := db.First(&regDetails).Error; err != nil {
		return nil, fmt.Errorf("failed to load registration details: %v", err)
	}

	capacities := map[int]int{
		MESS_A_LDH: regDetails.MessALDHCapacity,
		MESS_A_UDH: regDetails.MessAUDHCapacity,
		MESS_B_LDH: regDetails.MessBLDHCapacity,
		MESS_B_UDH: regDetails.MessBUDHCapacity,
		VEG_MESS:   regDetails.MessALDHCapacity + regDetails.MessAUDHCapacity, // Veg mess shares capacity
	}

	return capacities, nil
}

// InitializeMessCapacities sets the capacity limits for each mess in Redis
func (r *RedisMessService) InitializeMessCapacities(capacities map[int]int) error {
	pipe := r.client.Pipeline()

	for messID, capacity := range capacities {
		key := fmt.Sprintf(MESS_CAPACITY_KEY, messID)
		pipe.Set(r.ctx, key, capacity, 0)
	}

	_, err := pipe.Exec(r.ctx)
	return err
}

// RefreshCapacitiesFromDB loads latest capacities from database and updates Redis
func (r *RedisMessService) RefreshCapacitiesFromDB() error {
	capacities, err := r.LoadCapacitiesFromDB()
	if err != nil {
		return err
	}

	return r.InitializeMessCapacities(capacities)
}

// AttemptMessRegistration tries to register a user for a mess using atomic operations
func (r *RedisMessService) AttemptMessRegistration(userID uint, messID int) (bool, error) {
	// Validate mess ID first
	if !IsValidMessID(messID) {
		return false, fmt.Errorf("invalid mess ID: %d. Valid range is 1-5", messID)
	}

	counterKey := fmt.Sprintf(MESS_COUNTER_KEY, messID)
	capacityKey := fmt.Sprintf(MESS_CAPACITY_KEY, messID)
	userMessKey := fmt.Sprintf(USER_MESS_KEY, userID)
	messName := GetMessName(messID)

	// Use Redis transaction to ensure atomicity
	txf := func(tx *redis.Tx) error {
		// Check if user already has a mess assigned
		existingMess, err := tx.Get(r.ctx, userMessKey).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if existingMess != "" && existingMess != "0" {
			existingMessID, _ := strconv.Atoi(existingMess)
			return fmt.Errorf("user already has %s assigned", GetMessName(existingMessID))
		}

		// Get current count and capacity
		currentCount, err := tx.Get(r.ctx, counterKey).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		if err == redis.Nil {
			currentCount = 0
		}

		capacity, err := tx.Get(r.ctx, capacityKey).Int()
		if err != nil {
			return fmt.Errorf("capacity not initialized for %s", messName)
		}

		// Check if mess is full
		if currentCount >= capacity {
			return fmt.Errorf("%s is full (%d/%d)", messName, currentCount, capacity)
		}

		// If we reach here, registration is possible
		// Use MULTI/EXEC to perform atomic updates
		_, err = tx.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
			// Increment mess counter
			pipe.Incr(r.ctx, counterKey)
			// Set user's mess
			pipe.Set(r.ctx, userMessKey, messID, 0)
			// Add user to pending sync queue (as string for consistency)
			pipe.SAdd(r.ctx, PENDING_SYNC_KEY, strconv.FormatUint(uint64(userID), 10))
			return nil
		})

		return err
	}

	// Retry the transaction up to 3 times
	for i := 0; i < 3; i++ {
		err := r.client.Watch(r.ctx, txf, userMessKey, counterKey)
		if err == nil {
			return true, nil // Registration successful
		}
		if err == redis.TxFailedErr {
			// Transaction failed due to watched keys being modified, retry
			continue
		}
		// Other errors are not retryable
		return false, err
	}

	return false, fmt.Errorf("registration failed after retries")
}

// GetUserMess retrieves the mess assignment for a user from Redis
func (r *RedisMessService) GetUserMess(userID uint) (int, error) {
	userMessKey := fmt.Sprintf(USER_MESS_KEY, userID)
	result, err := r.client.Get(r.ctx, userMessKey).Result()

	if err == redis.Nil {
		return 0, nil // User has no mess assigned
	}
	if err != nil {
		return 0, err
	}

	messID, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("invalid mess ID format: %s", result)
	}

	return messID, nil
}

// GetMessCount returns the current registration count for a mess
func (r *RedisMessService) GetMessCount(messID int) (int, error) {
	counterKey := fmt.Sprintf(MESS_COUNTER_KEY, messID)
	count, err := r.client.Get(r.ctx, counterKey).Int()

	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

// GetMessCapacity returns the capacity limit for a mess
func (r *RedisMessService) GetMessCapacity(messID int) (int, error) {
	capacityKey := fmt.Sprintf(MESS_CAPACITY_KEY, messID)
	capacity, err := r.client.Get(r.ctx, capacityKey).Int()

	if err == redis.Nil {
		return 0, fmt.Errorf("capacity not set for mess %d", messID)
	}
	return capacity, err
}

// GetPendingSyncUsers returns all user IDs that need to be synced to the database
func (r *RedisMessService) GetPendingSyncUsers() ([]string, error) {
	return r.client.SMembers(r.ctx, PENDING_SYNC_KEY).Result()
}

// RemoveFromPendingSync removes a user from the pending sync queue
func (r *RedisMessService) RemoveFromPendingSync(userID uint) error {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	result := r.client.SRem(r.ctx, PENDING_SYNC_KEY, userIDStr)
	if result.Err() != nil {
		return result.Err()
	}

	// Manually remove key
	userMessKey := fmt.Sprintf(USER_MESS_KEY, userID)
	res := r.client.Del(r.ctx, userMessKey)
	if res.Err() != nil {
		return res.Err()
	}

	// Log if the user was actually removed (result should be 1 if removed, 0 if not found)
	if result.Val() == 0 {
		return fmt.Errorf("user %d not found in pending sync queue", userID)
	}

	return nil

}

// ClearUserRegistration removes a user's mess assignment (for rollback scenarios)
func (r *RedisMessService) ClearUserRegistration(userID uint, messID int) error {
	counterKey := fmt.Sprintf(MESS_COUNTER_KEY, messID)
	userMessKey := fmt.Sprintf(USER_MESS_KEY, userID)

	pipe := r.client.Pipeline()
	pipe.Decr(r.ctx, counterKey)
	pipe.Del(r.ctx, userMessKey)
	pipe.SRem(r.ctx, PENDING_SYNC_KEY, strconv.FormatUint(uint64(userID), 10))

	_, err := pipe.Exec(r.ctx)
	return err
}

// MessStats represents statistics for a single mess
type MessStats struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Capacity  int    `json:"capacity"`
	Available int    `json:"available"`
	FullPct   int    `json:"full_percentage"`
}

// GetAllMessStats returns registration statistics for all messes
func (r *RedisMessService) GetAllMessStats() (map[string]MessStats, error) {
	stats := make(map[string]MessStats)

	for messID := MESS_A_LDH; messID <= MESS_B_UDH; messID++ {
		count, err := r.GetMessCount(messID)
		if err != nil {
			return nil, err
		}

		capacity, err := r.GetMessCapacity(messID)
		if err != nil {
			return nil, err
		}

		available := capacity - count
		fullPct := 0
		if capacity > 0 {
			fullPct = (count * 100) / capacity
		}

		messName := GetMessName(messID)
		stats[messName] = MessStats{
			ID:        messID,
			Name:      messName,
			Count:     count,
			Capacity:  capacity,
			Available: available,
			FullPct:   fullPct,
		}
	}

	return stats, nil
}

// GetMessStatsByGroup returns statistics grouped by mess (A/B) with LDH/UDH breakdown
func (r *RedisMessService) GetMessStatsByGroup() (map[string]map[string]MessStats, error) {
	allStats, err := r.GetAllMessStats()
	if err != nil {
		return nil, err
	}

	grouped := map[string]map[string]MessStats{
		"MessA": make(map[string]MessStats),
		"MessB": make(map[string]MessStats),
	}

	for _, stat := range allStats {
		switch stat.ID {
		case MESS_A_LDH:
			grouped["MessA"]["LDH"] = stat
		case MESS_A_UDH:
			grouped["MessA"]["UDH"] = stat
		case MESS_B_LDH:
			grouped["MessB"]["LDH"] = stat
		case MESS_B_UDH:
			grouped["MessB"]["UDH"] = stat
		}
	}

	return grouped, nil
}
