package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateRollNumber validates the format of a roll number
func ValidateRollNumber(rollNo string) error {
	if rollNo == "" {
		return fmt.Errorf("roll number cannot be empty")
	}

	// Example: ES23BTECH1028 (adjust pattern based on your college format)
	pattern := `^[A-Z]{2}[0-9]{2}[A-Z]{1,6}[0-9]{4}$`
	matched, err := regexp.MatchString(pattern, strings.ToUpper(rollNo))
	if err != nil {
		return fmt.Errorf("error validating roll number: %v", err)
	}

	if !matched {
		return fmt.Errorf("invalid roll number format")
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return fmt.Errorf("error validating email: %v", err)
	}

	if !matched {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateVegType validates veg/non-veg type
func ValidateVegType(vegType string) error {
	validTypes := []string{"veg", "non-veg"}
	for _, valid := range validTypes {
		if strings.ToLower(vegType) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid veg type, must be 'veg' or 'non-veg'")
}

// ValidateUserType validates user type
func ValidateUserType(userType string) error {
	validTypes := []string{"student", "admin", "mess_staff"}
	for _, valid := range validTypes {
		if strings.ToLower(userType) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid user type, must be 'student', 'admin', or 'mess_staff'")
}

// ValidateMessNumber validates mess number
func ValidateMessNumber(mess int) error {
	if mess < 0 || mess > 3 {
		return fmt.Errorf("invalid mess number, must be between 0 and 3")
	}
	return nil
}

// ValidateSwapType validates swap request type
func ValidateSwapType(swapType string) error {
	validTypes := []string{"friend", "public"}
	for _, valid := range validTypes {
		if strings.ToLower(swapType) == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid swap type, must be 'friend' or 'public'")
}

// SanitizeInput removes extra spaces and converts to proper case
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// SanitizeName converts name to proper case
func SanitizeName(name string) string {
	name = SanitizeInput(name)
	words := strings.Fields(name)
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}

// SanitizeRollNumber converts roll number to uppercase
func SanitizeRollNumber(rollNo string) string {
	return strings.ToUpper(SanitizeInput(rollNo))
}

// SanitizeEmail converts email to lowercase
func SanitizeEmail(email string) string {
	return strings.ToLower(SanitizeInput(email))
}
