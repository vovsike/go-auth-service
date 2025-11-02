package apperror

import (
	"errors"
	"net/http"
	"testing"
)

// sentinel error for testing
var errSentinel = errors.New("sentinel")

func TestHTTPError_Error_PrefersMessage(t *testing.T) {
	he := &HTTPError{Err: errSentinel, StatusCode: http.StatusBadRequest, Message: "custom message"}
	if he.Error() != "custom message" {
		t.Fatalf("Error() should return Message when set; got %q", he.Error())
	}
}

func TestHTTPError_Error_FallsBackToErr(t *testing.T) {
	he := &HTTPError{Err: errSentinel, StatusCode: http.StatusUnauthorized}
	if he.Error() != errSentinel.Error() {
		t.Fatalf("Error() should return underlying error message when Message is empty; got %q", he.Error())
	}
}

func TestHTTPError_Unwrap_AllowsErrorsIsAndAs(t *testing.T) {
	he := &HTTPError{Err: errSentinel, StatusCode: http.StatusForbidden}

	if !errors.Is(he, errSentinel) {
		t.Fatalf("errors.Is should match the wrapped error")
	}

	var target *HTTPError
	if !errors.As(he, &target) {
		t.Fatalf("errors.As should unwrap into *HTTPError")
	}
	if target.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, target.StatusCode)
	}
}

func TestNewHTTPError_SetsFields(t *testing.T) {
	he := NewHTTPError(errSentinel, http.StatusInternalServerError)
	if he.Err != errSentinel {
		t.Fatalf("Err not set correctly")
	}
	if he.StatusCode != http.StatusInternalServerError {
		t.Fatalf("StatusCode not set correctly; got %d", he.StatusCode)
	}
	if he.Message != "" {
		t.Fatalf("Message should be empty; got %q", he.Message)
	}
}

func TestNewHTTPErrorWithMessage_SetsFields(t *testing.T) {
	he := NewHTTPErrorWithMessage(errSentinel, http.StatusTeapot, "i am a teapot")
	if he.Err != errSentinel {
		t.Fatalf("Err not set correctly")
	}
	if he.StatusCode != http.StatusTeapot {
		t.Fatalf("StatusCode not set correctly; got %d", he.StatusCode)
	}
	if he.Message != "i am a teapot" {
		t.Fatalf("Message not set correctly; got %q", he.Message)
	}
}

func TestHelperConstructors_SetStatusCodesAndPropagateError(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(error) *HTTPError
		wantStatus int
	}{
		{name: "BadRequest", fn: BadRequest, wantStatus: http.StatusBadRequest},
		{name: "Unauthorized", fn: Unauthorized, wantStatus: http.StatusUnauthorized},
		{name: "Forbidden", fn: Forbidden, wantStatus: http.StatusForbidden},
		{name: "NotFound", fn: NotFound, wantStatus: http.StatusNotFound},
		{name: "InternalServerError", fn: InternalServerError, wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			he := tt.fn(errSentinel)
			if he.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, he.StatusCode)
			}
			if !errors.Is(he, errSentinel) {
				t.Fatalf("wrapped error not propagated correctly")
			}
		})
	}
}
