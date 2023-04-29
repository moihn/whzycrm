package util

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	Code    int
	Message string `json:"message"`
	ErrRef  string `json:"reference,omitempty"`
}

func NewInternalServiceErrorStatus(message string, ref string) Status {
	return Status{
		Code:    http.StatusInternalServerError,
		Message: message,
		ErrRef:  ref,
	}
}

func NewServiceUnavaulableStatus(message string, ref string) Status {
	return Status{
		Code:    http.StatusServiceUnavailable,
		Message: message,
		ErrRef:  ref,
	}
}

func NewBadRequestStatus(message string) Status {
	return Status{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundStatus(message string) Status {
	return Status{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func (status Status) ErrorBytes() []byte {
	jsonByte, _ := json.Marshal(status)
	return jsonByte
}

func (status Status) Error() string {
	return string(status.ErrorBytes())
}

func (status *Status) ErrorRef() string {
	return status.ErrRef
}
