package db

import (
	"time"

	"gorm.io/gorm"
)

// Scans model
type Scans struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	MessID    uint      `json:"mess_id" gorm:"not null;check:mess_id >= 1 AND mess_id <= 4"`
	Meal      int       `json:"meal" gorm:"not null;check:meal >= 1 AND meal <= 4"`
	Date      time.Time `json:"date" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// LogCurrentMeal logs a scan for the user based on the current time
func LogCurrentMeal(db *gorm.DB, userID uint, messID uint) (*Scans, error) {

	istLocation := time.FixedZone("IST", 5*60*60+30*60)
	now := time.Now().In(istLocation)
	meal := 0

	switch {
	case now.Hour() >= 8 && now.Hour() <= 10 && (now.Hour() != 10 || now.Minute() <= 30):
		meal = 1 // Breakfast
	case now.Hour() >= 12 && now.Hour() <= 14 && (now.Hour() != 14 || now.Minute() <= 45):
		meal = 2 // Lunch
	case now.Hour() >= 17 && now.Hour() <= 18:
		meal = 3 // Snacks
	case now.Hour() >= 19 && now.Hour() <= 21 && (now.Hour() != 21 || now.Minute() <= 30):
		meal = 4 // Dinner
	default:
		return nil, nil // Not a meal time
	}

	// Check if already logged
	var existing Scans
	date := now.Truncate(24 * time.Hour) // only date part
	err := db.Where("user_id = ? AND date = ? AND meal = ?", userID, date, meal).First(&existing).Error
	if err == nil {
		// Already logged
		return &existing, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Insert new scan
	newScan := Scans{
		UserID: userID,
		MessID: messID,
		Meal:   meal,
		Date:   date,
	}
	if err := db.Create(&newScan).Error; err != nil {
		return nil, err
	}

	return &newScan, nil
}
