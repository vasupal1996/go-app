package auth

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

// HeaderRequestID header name to look for request id for request tracking
const HeaderRequestID = "X-Request-ID"

// TokenAuth defines method for implementing token authentication
type TokenAuth interface {
	SignToken(*jwt.Token) (string, error)
	VerifyToken(string) error
	GetClaim() interface{}
}

// Claim defines custom token claim type methods.
// Note: this claim is used to automatically parse token into struct when a request has jwt token in header
type Claim interface {
	ToJSON() string
	GetJWTToken() *jwt.Token
	IsAdmin() bool
}

// JWTToken represents jwt encoded token string for json format
type JWTToken struct {
	Token string `json:"token"`
}

// ToJSON := converting struct to json
func (jwt *JWTToken) ToJSON() string {
	json, _ := json.Marshal(jwt.Token)
	return string(json)
}
