package utils

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ArchiveTable archives a table to a new table with format tablename_<month>_<year>
// This function copies the entire table structure and data to the new archived table
func ArchiveTable(db *gorm.DB, tableName string, archiveDate time.Time) (string, error) {
	// Format: tablename_November_2025
	month := archiveDate.Format("January")
	year := archiveDate.Format("2006")
	archiveTableName := fmt.Sprintf("%s_%s_%s", tableName, month, year)

	// Check if archive table already exists
	type TableExists struct {
		Exists bool
	}
	var result TableExists

	if err := db.Table("information_schema.tables").
		Select("COUNT(*) > 0 as exists").
		Where("table_schema = ? AND table_name = ?", "public", archiveTableName).
		Scan(&result).Error; err != nil {
		return "", fmt.Errorf("error checking if archive table exists: %w", err)
	}

	if result.Exists {
		return archiveTableName, fmt.Errorf("archive table %s already exists", archiveTableName)
	}

	// Create archive table as a copy of the original table structure and data
	createQuery := fmt.Sprintf("CREATE TABLE %s AS TABLE %s", archiveTableName, tableName)
	if err := db.Exec(createQuery).Error; err != nil {
		return "", fmt.Errorf("error creating archive table: %w", err)
	}

	return archiveTableName, nil
}

// ArchiveUsersTable archives the users table
func ArchiveUsersTable(db *gorm.DB, archiveDate time.Time) (string, error) {
	return ArchiveTable(db, "users", archiveDate)
}

// ArchiveScansTable archives the scans table
func ArchiveScansTable(db *gorm.DB, archiveDate time.Time) (string, error) {
	return ArchiveTable(db, "scans", archiveDate)
}

// ClearTableData clears all data from a table without dropping it
func ClearTableData(db *gorm.DB, tableName string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)
	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("error clearing table data: %w", err)
	}
	return nil
}

// ListArchivedTables lists all archived tables for a given base table name
func ListArchivedTables(db *gorm.DB, baseTableName string) ([]string, error) {
	var tables []string
	pattern := fmt.Sprintf("%s_%%", baseTableName)

	type TableInfo struct {
		TableName string `gorm:"column:table_name"`
	}
	var tableInfos []TableInfo

	if err := db.Table("information_schema.tables").
		Select("table_name").
		Where("table_schema = ? AND table_name LIKE ?", "public", pattern).
		Order("table_name DESC").
		Find(&tableInfos).Error; err != nil {
		return nil, fmt.Errorf("error listing archived tables: %w", err)
	}

	for _, ti := range tableInfos {
		tables = append(tables, ti.TableName)
	}

	return tables, nil
}

// GetArchivedTableData retrieves data from an archived table
func GetArchivedTableData(db *gorm.DB, archiveTableName string, dest interface{}) error {
	if err := db.Table(archiveTableName).Find(dest).Error; err != nil {
		return fmt.Errorf("error retrieving archived data: %w", err)
	}
	return nil
}

// ArchiveCycleData archives both users and scans data for a cycle
// Returns the names of the archived tables
func ArchiveCycleData(db *gorm.DB, archiveDate time.Time) (map[string]string, error) {
	result := make(map[string]string)

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	// Archive scans table
	scansTableName, err := ArchiveScansTable(tx, archiveDate)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error archiving scans table: %w", err)
	}
	result["scans"] = scansTableName

	// Archive users table (creates a snapshot)
	usersTableName, err := ArchiveUsersTable(tx, archiveDate)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error archiving users table: %w", err)
	}
	result["users"] = usersTableName

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return result, nil
}

// PrepareNewCycle archives old data and prepares for a new registration cycle
// This clears scans and resets next_mess in users
func PrepareNewCycle(db *gorm.DB, archiveDate time.Time) (map[string]string, error) {
	// Archive the cycle data
	archivedTables, err := ArchiveCycleData(db, archiveDate)
	if err != nil {
		return nil, fmt.Errorf("error archiving cycle data: %w", err)
	}

	// Start a transaction for cleanup
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("error starting cleanup transaction: %w", tx.Error)
	}

	// Clear scans table
	if err := ClearTableData(tx, "scans"); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error clearing scans: %w", err)
	}

	// Update users: set mess = next_mess and reset next_mess to 0
	if err := tx.Model(&struct {
		Mess     int8 `gorm:"column:mess"`
		NextMess int8 `gorm:"column:next_mess"`
	}{}).
		Table("users").
		Where("next_mess != ?", 0).
		Updates(map[string]interface{}{
			"mess":      gorm.Expr("next_mess"),
			"next_mess": 0,
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating user mess assignments: %w", err)
	}

	// Commit cleanup transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing cleanup transaction: %w", err)
	}

	return archivedTables, nil
}
