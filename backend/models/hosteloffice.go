package models

type UserInfo struct {
	Name   string
	Email  string
	RollNo string
	Mess   int8
}

type EditUserInfo struct {
	RollNo      string `json:"roll_no"`      //To identify the student
	Mess        int8   `json:"mess"`         //The hostel office can change this value
	CanRegister bool   `json:"can_register"` //fasle -> The user has been deactivated.
}

// DateRangeRequest represents a date range query
type DateRangeRequest struct {
	FromDate string `json:"from_date" form:"from_date" binding:"required"` // Format: 2006-01-02
	ToDate   string `json:"to_date" form:"to_date" binding:"required"`     // Format: 2006-01-02
}

// ArchiveRequest represents a request to archive cycle data
type ArchiveRequest struct {
	Month int `json:"month" binding:"required,min=1,max=12"` // 1-12
	Year  int `json:"year" binding:"required,min=2020"`      // Year to archive
}

// BulkEditRequest represents a request to bulk edit users
type BulkEditRequest struct {
	Updates []EditUserInfo `json:"updates" binding:"required"`
}

// UploadResponse represents the response after uploading a CSV file
type UploadResponse struct {
	Message       string   `json:"message"`
	RecordsAdded  int      `json:"records_added,omitempty"`
	RecordsUpdate int      `json:"records_updated,omitempty"`
	Errors        []string `json:"errors,omitempty"`
}

// ArchiveResponse represents the response after archiving data
type ArchiveResponse struct {
	Message        string            `json:"message"`
	ArchivedTables map[string]string `json:"archived_tables"`
}

// ArchivedTablesListResponse represents the list of archived tables
type ArchivedTablesListResponse struct {
	Users []string `json:"users"`
	Scans []string `json:"scans"`
}
