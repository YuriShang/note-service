package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound  = NewAppError("not found", "")
	ErrForbidden = NewAppError("forbidden", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

func NewAppError(message, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, "some thing wrong with user data")
}

func NotFoundError(message string) *AppError {
	return NewAppError(message, "")
}

func systemError(developerMessage string) *AppError {
	return NewAppError("system error", developerMessage)
}

func APIError(message, developerMessage string) *AppError {
	return NewAppError(message, developerMessage)
}
