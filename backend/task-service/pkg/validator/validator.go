package validator

import (
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

func IsNotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

func HasMinLength(value string, min int) bool {
	return len(value) >= min
}

func HasMaxLength(value string, max int) bool {
	return len(value) <= max
}

func IsValidPriority(priority string) bool {
	validPriorities := []string{"low", "medium", "high", "urgent"}
	for _, p := range validPriorities {
		if p == priority {
			return true
		}
	}
	return false
}

func IsValidStatus(status string) bool {
	validStatuses := []string{"todo", "in_progress", "review", "done"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
