package db

type LoggerDetails struct {
	UserID      uint        `json:"user_id"`
	Action      string      `json:"action"`
	Message     string      `json:"message"`
	IPAddress   string      `json:"ip_address"`
	HTTPDetails HTTPDetails `gorm:"embedded;embeddedPrefix:http_"`
	Timestamp   string      `json:"timestamp"`
}

type HTTPDetails struct {
	Method     string `json:"method"`
	Endpoint   string `json:"endpoint"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
