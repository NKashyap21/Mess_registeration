package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
)

// ParseStudentsCSV parses a CSV file containing student data
// Expected format: Name,Email,Phone,RollNo,Mess,Type,CanRegister
func ParseStudentsCSV(file multipart.File) ([]models.User, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Skip header row
	var users []models.User
	for i, record := range records[1:] {
		if len(record) < 7 {
			return nil, fmt.Errorf("row %d has insufficient columns", i+2)
		}

		mess, err := strconv.ParseInt(strings.TrimSpace(record[4]), 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid mess value at row %d: %w", i+2, err)
		}

		userType, err := strconv.ParseInt(strings.TrimSpace(record[5]), 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid type value at row %d: %w", i+2, err)
		}

		canRegister, err := strconv.ParseBool(strings.TrimSpace(record[6]))
		if err != nil {
			return nil, fmt.Errorf("invalid can_register value at row %d: %w", i+2, err)
		}

		phone := strings.TrimSpace(record[2])
		user := models.User{
			Name:        strings.TrimSpace(record[0]),
			Email:       strings.TrimSpace(record[1]),
			Phone:       &phone,
			RollNo:      strings.TrimSpace(record[3]),
			Mess:        int8(mess),
			Type:        int8(userType),
			CanRegister: canRegister,
		}
		users = append(users, user)
	}

	return users, nil
}

// ParseCanRegisterCSV parses a CSV file for updating can_register status
// Expected format: RollNo,CanRegister
func ParseCanRegisterCSV(file multipart.File) (map[string]bool, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Skip header row
	updates := make(map[string]bool)
	for i, record := range records[1:] {
		if len(record) < 2 {
			return nil, fmt.Errorf("row %d has insufficient columns", i+2)
		}

		canRegister, err := strconv.ParseBool(strings.TrimSpace(record[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid can_register value at row %d: %w", i+2, err)
		}

		rollNo := strings.TrimSpace(record[0])
		updates[rollNo] = canRegister
	}

	return updates, nil
}

// ExportUsersToCSV exports user data to CSV format in memory
func ExportUsersToCSV(users []models.User) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	// Write header
	header := []string{"ID", "Name", "Email", "Phone", "RollNo", "Mess", "NextMess", "Type", "CanRegister", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data rows
	for _, user := range users {
		phone := ""
		if user.Phone != nil {
			phone = *user.Phone
		}

		record := []string{
			strconv.FormatUint(uint64(user.ID), 10),
			user.Name,
			user.Email,
			phone,
			user.RollNo,
			strconv.FormatInt(int64(user.Mess), 10),
			strconv.FormatInt(int64(user.NextMess), 10),
			strconv.FormatInt(int64(user.Type), 10),
			strconv.FormatBool(user.CanRegister),
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return buffer, nil
}

// ScanRecord represents a scan record for CSV export
type ScanRecord struct {
	ID        uint
	UserID    uint
	UserName  string
	RollNo    string
	MessID    uint
	Meal      int
	Date      time.Time
	CreatedAt time.Time
}

// ExportScansToCSV exports scan data to CSV format in memory
func ExportScansToCSV(scans []ScanRecord) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	// Write header
	header := []string{"ID", "UserID", "UserName", "RollNo", "MessID", "Meal", "Date", "CreatedAt"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data rows
	for _, scan := range scans {
		record := []string{
			strconv.FormatUint(uint64(scan.ID), 10),
			strconv.FormatUint(uint64(scan.UserID), 10),
			scan.UserName,
			scan.RollNo,
			strconv.FormatUint(uint64(scan.MessID), 10),
			strconv.Itoa(scan.Meal),
			scan.Date.Format("2006-01-02"),
			scan.CreatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return buffer, nil
}

// RegistrationRecord represents a registration record for CSV export
type RegistrationRecord struct {
	UserID    uint
	UserName  string
	RollNo    string
	Email     string
	Mess      int8
	NextMess  int8
	UpdatedAt time.Time
}

// ExportRegistrationsToCSV exports registration data to CSV format in memory
func ExportRegistrationsToCSV(registrations []RegistrationRecord) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	// Write header
	header := []string{"UserID", "UserName", "RollNo", "Email", "CurrentMess", "RegisteredMess", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data rows
	for _, reg := range registrations {
		record := []string{
			strconv.FormatUint(uint64(reg.UserID), 10),
			reg.UserName,
			reg.RollNo,
			reg.Email,
			strconv.FormatInt(int64(reg.Mess), 10),
			strconv.FormatInt(int64(reg.NextMess), 10),
			reg.UpdatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("error writing CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return buffer, nil
}
