package main

import (
	"awesomeProject/internal/apperror"
	"errors"
	"net/http"
)

func ErrorHandler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var httpErr *apperror.HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Error(), httpErr.StatusCode)
				return
			}
			// Default to 500 for unknown errors
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
