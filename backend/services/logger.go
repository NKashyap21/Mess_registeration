package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoggerService struct {
	db          *gorm.DB
	logChan     chan models.LoggerDetails
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	batchSize   int
	flushTime   time.Duration
	buffer      []models.LoggerDetails
	bufferMutex sync.Mutex
}

var (
	loggerInstance *LoggerService
	loggerOnce     sync.Once
)

// GetLoggerService returns a singleton instance of the logger service
func GetLoggerService() *LoggerService {
	loggerOnce.Do(func() {
		loggerInstance = NewLoggerService()
	})
	return loggerInstance
}

// NewLoggerService creates a new logger service instance
func NewLoggerService() *LoggerService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &LoggerService{
		db:        config.GetDB(),
		logChan:   make(chan models.LoggerDetails, 1000), // Buffer for 1000 log entries
		ctx:       ctx,
		cancel:    cancel,
		batchSize: 50,              // Batch size for bulk inserts
		flushTime: 5 * time.Second, // Flush every 5 seconds
		buffer:    make([]models.LoggerDetails, 0),
	}

	// Start background worker
	service.startBackgroundWorker()

	return service
}

// startBackgroundWorker starts the background goroutine that processes logs
func (ls *LoggerService) startBackgroundWorker() {
	ls.wg.Add(1)
	go func() {
		defer ls.wg.Done()

		ticker := time.NewTicker(ls.flushTime)
		defer ticker.Stop()

		for {
			select {
			case logEntry := <-ls.logChan:
				ls.addToBuffer(logEntry)

				// Flush if buffer is full
				ls.bufferMutex.Lock()
				if len(ls.buffer) >= ls.batchSize {
					ls.flushBuffer()
				}
				ls.bufferMutex.Unlock()

			case <-ticker.C:
				// Periodic flush
				ls.bufferMutex.Lock()
				if len(ls.buffer) > 0 {
					ls.flushBuffer()
				}
				ls.bufferMutex.Unlock()

			case <-ls.ctx.Done():
				// Final flush before shutdown
				ls.bufferMutex.Lock()
				if len(ls.buffer) > 0 {
					ls.flushBuffer()
				}
				ls.bufferMutex.Unlock()
				return
			}
		}
	}()
}

// addToBuffer adds a log entry to the buffer
func (ls *LoggerService) addToBuffer(logEntry models.LoggerDetails) {
	ls.bufferMutex.Lock()
	defer ls.bufferMutex.Unlock()
	ls.buffer = append(ls.buffer, logEntry)
}

// flushBuffer writes all buffered logs to database
func (ls *LoggerService) flushBuffer() {
	if len(ls.buffer) == 0 {
		return
	}

	// Make a copy of buffer and clear it
	logs := make([]models.LoggerDetails, len(ls.buffer))
	copy(logs, ls.buffer)
	ls.buffer = ls.buffer[:0] // Clear buffer

	// Write to database in batch
	if err := ls.db.CreateInBatches(logs, len(logs)).Error; err != nil {
		log.Printf("Failed to write logs to database: %v", err)
		// In case of error, we could implement retry logic or write to file
	} else {
		log.Printf("Successfully wrote %d log entries to database", len(logs))
	}
}

// LogHTTPRequest logs HTTP request details
func (ls *LoggerService) LogHTTPRequest(c *gin.Context, userID uint, statusCode int, message string) {
	logEntry := models.LoggerDetails{
		UserID:    userID,
		Action:    "HTTP_REQUEST",
		Message:   message,
		IPAddress: c.ClientIP(),
		HTTPDetails: models.HTTPDetails{
			Method:     c.Request.Method,
			Endpoint:   c.Request.URL.Path,
			StatusCode: statusCode,
			Message:    message,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Send to channel (non-blocking)
	select {
	case ls.logChan <- logEntry:
	default:
		log.Println("Warning: Log channel is full, dropping log entry")
	}
}

// LogUserAction logs user-specific actions
func (ls *LoggerService) LogUserAction(userID uint, action, message, ipAddress string) {
	logEntry := models.LoggerDetails{
		UserID:    userID,
		Action:    action,
		Message:   message,
		IPAddress: ipAddress,
		HTTPDetails: models.HTTPDetails{
			Method:     "",
			Endpoint:   "",
			StatusCode: 0,
			Message:    message,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	select {
	case ls.logChan <- logEntry:
	default:
		log.Println("Warning: Log channel is full, dropping log entry")
	}
}

// LogSystemAction logs system-level actions
func (ls *LoggerService) LogSystemAction(action, message string) {
	logEntry := models.LoggerDetails{
		UserID:    0, // System actions have no user ID
		Action:    action,
		Message:   message,
		IPAddress: "system",
		HTTPDetails: models.HTTPDetails{
			Method:     "",
			Endpoint:   "",
			StatusCode: 0,
			Message:    message,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	select {
	case ls.logChan <- logEntry:
	default:
		log.Println("Warning: Log channel is full, dropping log entry")
	}
}

// LogDatabaseAction logs database operations
func (ls *LoggerService) LogDatabaseAction(userID uint, action, table, message, ipAddress string) {
	logEntry := models.LoggerDetails{
		UserID:    userID,
		Action:    fmt.Sprintf("DB_%s_%s", action, table),
		Message:   message,
		IPAddress: ipAddress,
		HTTPDetails: models.HTTPDetails{
			Method:     "",
			Endpoint:   "",
			StatusCode: 0,
			Message:    message,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	select {
	case ls.logChan <- logEntry:
	default:
		log.Println("Warning: Log channel is full, dropping log entry")
	}
}

// LogAuthAction logs authentication-related actions
func (ls *LoggerService) LogAuthAction(userID uint, action, message, ipAddress string) {
	logEntry := models.LoggerDetails{
		UserID:    userID,
		Action:    fmt.Sprintf("AUTH_%s", action),
		Message:   message,
		IPAddress: ipAddress,
		HTTPDetails: models.HTTPDetails{
			Method:     "",
			Endpoint:   "",
			StatusCode: 0,
			Message:    message,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	select {
	case ls.logChan <- logEntry:
	default:
		log.Println("Warning: Log channel is full, dropping log entry")
	}
}

// GetLogs retrieves logs from database with pagination and filtering
func (ls *LoggerService) GetLogs(limit, offset int, userID *uint, action string, startDate, endDate *time.Time) ([]models.LoggerDetails, int64, error) {
	var logs []models.LoggerDetails
	var total int64

	query := ls.db.Model(&models.LoggerDetails{})

	// Apply filters
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if action != "" {
		query = query.Where("action ILIKE ?", "%"+action+"%")
	}

	if startDate != nil {
		query = query.Where("timestamp >= ?", startDate.Format(time.RFC3339))
	}

	if endDate != nil {
		query = query.Where("timestamp <= ?", endDate.Format(time.RFC3339))
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogsByDateRange retrieves logs within a specific date range
func (ls *LoggerService) GetLogsByDateRange(startDate, endDate time.Time) ([]models.LoggerDetails, error) {
	var logs []models.LoggerDetails

	err := ls.db.Where("timestamp >= ? AND timestamp <= ?",
		startDate.Format(time.RFC3339),
		endDate.Format(time.RFC3339)).
		Order("timestamp DESC").
		Find(&logs).Error

	return logs, err
}

// GetUserActivity retrieves activity logs for a specific user
func (ls *LoggerService) GetUserActivity(userID uint, limit int) ([]models.LoggerDetails, error) {
	var logs []models.LoggerDetails

	err := ls.db.Where("user_id = ?", userID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// GetSystemLogs retrieves system-level logs
func (ls *LoggerService) GetSystemLogs(limit int) ([]models.LoggerDetails, error) {
	var logs []models.LoggerDetails

	err := ls.db.Where("user_id = 0").
		Order("timestamp DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// ExportLogs exports logs to JSON format
func (ls *LoggerService) ExportLogs(startDate, endDate time.Time) ([]byte, error) {
	logs, err := ls.GetLogsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(logs, "", "  ")
}

// GetLogStats returns statistics about logs
func (ls *LoggerService) GetLogStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total logs count
	var totalLogs int64
	if err := ls.db.Model(&models.LoggerDetails{}).Count(&totalLogs).Error; err != nil {
		return nil, err
	}
	stats["total_logs"] = totalLogs

	// Logs by action
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	if err := ls.db.Model(&models.LoggerDetails{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Order("count DESC").
		Limit(10).
		Find(&actionStats).Error; err != nil {
		return nil, err
	}
	stats["top_actions"] = actionStats

	// Logs by date (last 7 days)
	var dateStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	if err := ls.db.Model(&models.LoggerDetails{}).
		Select("DATE(timestamp) as date, COUNT(*) as count").
		Where("DATE(timestamp) >= ?", sevenDaysAgo).
		Group("DATE(timestamp)").
		Order("date DESC").
		Find(&dateStats).Error; err != nil {
		return nil, err
	}
	stats["daily_logs"] = dateStats

	return stats, nil
}

// Shutdown gracefully shuts down the logger service
func (ls *LoggerService) Shutdown() {
	log.Println("Shutting down logger service...")
	ls.cancel()
	ls.wg.Wait()
	close(ls.logChan)
	log.Println("Logger service shutdown complete")
}
