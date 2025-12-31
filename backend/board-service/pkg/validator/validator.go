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

func IsValidSlug(slug string) bool {
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return slugRegex.MatchString(slug)
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
