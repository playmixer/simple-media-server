package apperror

import (
	"encoding/json"
)

var (
	ErrPathNotFound       = NewAppError(nil, "path not found", "API-000001")
	ErrCantMarshal        = NewAppError(nil, "error marshal to json", "API-000002")
	ErrCantParseUrlParams = NewAppError(nil, "cant parse url params", "API-000003")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func NewAppError(err error, message, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}

}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
