package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"auth_service/constants"
)

// ValidatePasswordStrength checks password meets complexity requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < constants.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", constants.MinPasswordLength)
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character (!@#$%%^&*)")
	}

	return nil
}

// ValidateUsername checks username format and length
func ValidateUsername(username string) error {
	if len(username) < constants.MinUsernameLength {
		return fmt.Errorf("username must be at least %d characters long", constants.MinUsernameLength)
	}
	if len(username) > constants.MaxUsernameLength {
		return fmt.Errorf("username must be at most %d characters long", constants.MaxUsernameLength)
	}

	// Only alphanumeric and underscores
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}

// ValidateEndpointPath checks endpoint path is safe
func ValidateEndpointPath(path string) error {
	if len(path) < constants.MinPathLength {
		return fmt.Errorf("path must not be empty")
	}
	if len(path) > constants.MaxPathLength {
		return fmt.Errorf("path must be at most %d characters long", constants.MaxPathLength)
	}

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with /")
	}

	// Prevent path traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Prevent double slashes
	if strings.Contains(path, "//") {
		return fmt.Errorf("double slashes not allowed in path")
	}

	// Prevent null bytes
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("null bytes not allowed in path")
	}

	return nil
}

// ValidateEndpointTargets checks target URLs are valid
func ValidateEndpointTargets(targets []string) error {
	if len(targets) == 0 {
		return fmt.Errorf("at least one target endpoint required")
	}

	for i, target := range targets {
		if target == "" {
			return fmt.Errorf("target %d is empty", i+1)
		}
		if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
			return fmt.Errorf("target %d must start with http:// or https://", i+1)
		}
	}

	return nil
}

// ValidateHTTPMethod checks if method is valid
func ValidateHTTPMethod(method string) error {
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
		"ANY":     true,
	}

	if !validMethods[strings.ToUpper(method)] {
		return fmt.Errorf("invalid HTTP method: %s", method)
	}

	return nil
}
