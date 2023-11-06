// errors.go
package main

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Detail
}

func NewHTTPError(err error, status int, detail string) error {
	return &HTTPError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func handleError(w http.ResponseWriter, err error) {
	httpError, ok := err.(*HTTPError)
	if !ok {
		httpError = &HTTPError{
			Cause:  err,
			Detail: "Internal Server Error",
			Status: http.StatusInternalServerError,
		}
	}

	w.WriteHeader(httpError.Status)
	w.Header().Set("Content-Type", "application/json")

	errorResponse := map[string]interface{}{
		"error":   httpError.Detail,
		"message": httpError.Cause.Error(),
	}

	responseJSON, _ := json.Marshal(errorResponse)
	w.Write(responseJSON)
}
