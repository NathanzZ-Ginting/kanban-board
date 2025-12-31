package validator

import (
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, ", ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Email validation
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Username validation (alphanumeric, underscore, 3-30 chars)
func IsValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	return usernameRegex.MatchString(username)
}

// Password validation (min 8 chars)
func IsValidPassword(password string) bool {
	return len(password) >= 8
}

// Not empty validation
func IsNotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Min length validation
func HasMinLength(value string, min int) bool {
	return len(strings.TrimSpace(value)) >= min
}

// Max length validation
func HasMaxLength(value string, max int) bool {
	return len(strings.TrimSpace(value)) <= max
}
