package auth

import (
	"net/http"
)

// Session defines session storage methods
type Session interface {
	NewSessionID() string
	Create(string, http.ResponseWriter) error
	Get(*http.Request) (*TokenAuth, error)
	Update(http.ResponseWriter, *http.Request, *TokenAuth) error
	Delete(*http.Request) error
}
