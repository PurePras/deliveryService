package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
)

var (
	phonePattern   = regexp.MustCompile(`^[6-9]\d{9}$`)
	emailPattern   = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	otpCodePattern = regexp.MustCompile(`^\d{6}$`)
)

func requireUUID(field, value string) error {
	if _, err := uuid.Parse(value); err != nil {
		return fmt.Errorf("%w: %s must be a valid id", apperror.ErrValidation, field)
	}
	return nil
}

func requireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: %s is required", apperror.ErrValidation, field)
	}
	return nil
}

func requireNonNegativeNumber(field, value string) error {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil || n < 0 {
		return fmt.Errorf("%w: %s must be a non-negative number", apperror.ErrValidation, field)
	}
	return nil
}

func requirePositiveNumber(field, value string) error {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil || n <= 0 {
		return fmt.Errorf("%w: %s must be a positive number", apperror.ErrValidation, field)
	}
	return nil
}

func requireMinLength(field, value string, min int) error {
	if len(value) < min {
		return fmt.Errorf("%w: %s must be at least %d characters", apperror.ErrValidation, field, min)
	}
	return nil
}

// requirePhone validates a 10-digit Indian mobile number (starts 6-9).
func requirePhone(field, value string) error {
	if !phonePattern.MatchString(value) {
		return fmt.Errorf("%w: %s must be a valid 10-digit phone number", apperror.ErrValidation, field)
	}
	return nil
}

// requireEmail is only meant to be called when the field is present — callers with an
// optional email should skip calling this for a nil/empty value rather than requiring one.
func requireEmail(field, value string) error {
	if !emailPattern.MatchString(value) {
		return fmt.Errorf("%w: %s must be a valid email address", apperror.ErrValidation, field)
	}
	return nil
}

func requireOTPCode(field, value string) error {
	if !otpCodePattern.MatchString(value) {
		return fmt.Errorf("%w: %s must be a 6-digit code", apperror.ErrValidation, field)
	}
	return nil
}
