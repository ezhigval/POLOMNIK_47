package handlers

import (
	"errors"

	"polomnik/internal/application"
	"polomnik/internal/domain"
)

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func MapError(err error) *AppError {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &AppError{Status: 404, Code: "NOT_FOUND", Message: "Resource not found"}
	case errors.Is(err, application.ErrTourInactive):
		return &AppError{Status: 404, Code: "TOUR_INACTIVE", Message: "Tour is not active"}
	case errors.Is(err, application.ErrTourExpired):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Tour has already ended"}
	case errors.Is(err, application.ErrUnauthorized):
		return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	case errors.Is(err, domain.ErrInsufficientSlots):
		return &AppError{Status: 409, Code: "INSUFFICIENT_SLOTS", Message: "Not enough slots available"}
	case errors.Is(err, domain.ErrInvalidPeopleCount):
		return &AppError{Status: 422, Code: "INVALID_PEOPLE_COUNT", Message: "People count must be greater than 0"}
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Invalid booking status transition"}
	case errors.Is(err, domain.ErrInvalidBookingStatus):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Invalid booking status"}
	case errors.Is(err, domain.ErrInvalidCredentials):
		return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Invalid login or password"}
	case errors.Is(err, domain.ErrDuplicateEmail):
		return &AppError{Status: 409, Code: "DUPLICATE_EMAIL", Message: "Email is already registered"}
	case errors.Is(err, domain.ErrDuplicatePhone):
		return &AppError{Status: 409, Code: "DUPLICATE_PHONE", Message: "Phone is already registered"}
	case errors.Is(err, domain.ErrDuplicateSlug):
		return &AppError{Status: 409, Code: "DUPLICATE_SLUG", Message: "Page slug already exists"}
	case errors.Is(err, domain.ErrDuplicatePath):
		return &AppError{Status: 409, Code: "DUPLICATE_PATH", Message: "Page path already exists"}
	default:
		return mapValidationError(err)
	}
}

func MapBookingError(err error) *AppError {
	if errors.Is(err, domain.ErrNotFound) {
		return &AppError{Status: 404, Code: "TOUR_NOT_FOUND", Message: "Tour not found"}
	}
	return MapError(err)
}

func mapValidationError(err error) *AppError {
	switch {
	case errors.Is(err, domain.ErrInvalidID),
		errors.Is(err, domain.ErrInvalidSlug),
		errors.Is(err, domain.ErrInvalidTitle),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidCurrency),
		errors.Is(err, domain.ErrInvalidDateRange),
		errors.Is(err, domain.ErrInvalidSlots),
		errors.Is(err, domain.ErrInvalidContactName),
		errors.Is(err, domain.ErrInvalidPhone),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidPassword),
		errors.Is(err, domain.ErrInvalidTotalPrice),
		errors.Is(err, domain.ErrInvalidClientName),
		errors.Is(err, domain.ErrInvalidRating),
		errors.Is(err, domain.ErrInvalidReviewText),
		errors.Is(err, domain.ErrInvalidBlockType),
		errors.Is(err, domain.ErrInvalidPath):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: err.Error()}
	default:
		return nil
	}
}
