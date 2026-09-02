package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
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
