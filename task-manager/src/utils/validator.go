package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePassword checks if a password meets the minimum requirements
// Password requirements:
// - At least 8 characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, char := range password {
		if 'A' <= char && char <= 'Z' {
			hasUpper = true
		} else if 'a' <= char && char <= 'z' {
			hasLower = true
		} else if '0' <= char && char <= '9' {
			hasNumber = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}

	return true, ""
}

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove any HTML tags
	re := regexp.MustCompile("<[^>]*>")
	sanitized := re.ReplaceAllString(input, "")

	// Trim spaces
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// FormatValidationError formats validation errors from gin's binding into a user-friendly message
func FormatValidationError(err error) string {
	// Convert the error to string and clean it up
	errStr := err.Error()

	// Remove the common prefix from gin validation errors
	errStr = strings.Replace(errStr, "Key: ", "", -1)

	// Replace 'Error:' with ': '
	errStr = strings.Replace(errStr, "Error:", ":", -1)

	// Replace struct field names with more readable versions
	errStr = strings.Replace(errStr, "'ProjectInput.", "", -1)
	errStr = strings.Replace(errStr, "'UserRegistrationInput.", "", -1)
	errStr = strings.Replace(errStr, "'UserLoginInput.", "", -1)
	errStr = strings.Replace(errStr, "'", "", -1)

	return errStr
}
