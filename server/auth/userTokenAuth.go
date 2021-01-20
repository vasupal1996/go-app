package auth

import (
	"encoding/base64"
	"encoding/json"
	"go-app/server/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenAuthentication contains authentication related attributes and methods
type TokenAuthentication struct {
	Config *config.TokenAuthConfig
	User   *UserAuth
}

// NewTokenAuthentication returns new instance of TokenAuthentication
func NewTokenAuthentication(c *config.TokenAuthConfig) *TokenAuthentication {
	return &TokenAuthentication{Config: c}
}

// UserAuth contains encoded token info and user info
type UserAuth struct {
	UserClaim *UserClaim
	JWTToken  JWTToken
}

// UserClaim contains user related info for jwt token
type UserClaim struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	jwt.StandardClaims
}

// GetJWTToken return jwt.Token with claimInfo from user claim fields
func (uc *UserClaim) GetJWTToken() *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	return token
}

// ToJSON := converting struct to json
func (uc *UserClaim) ToJSON() string {
	json, _ := json.Marshal(uc)
	return string(json)
}

// IsAdmin if user is an admin user
func (uc *UserClaim) IsAdmin() bool {
	if uc.Type == "admin" {
		return true
	}
	return false
}

// SignToken sign and encodes jwt.Token as a string
func (t *TokenAuthentication) SignToken() (string, error) {
	if t.Config.JWTExpiresAt != 0 {
		expirationTime := time.Now().Add(time.Duration(t.Config.JWTExpiresAt) * time.Minute)
		t.User.UserClaim.StandardClaims.ExpiresAt = expirationTime.Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.User.UserClaim)
	tokenString, _ := token.SignedString([]byte(t.Config.JWTSignKey))
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), nil
}

// VerifyToken first verifies the authenticity of the jwt token string and then parse the token string into struct
func (t *TokenAuthentication) VerifyToken(tokenString string) error {
	uc := UserClaim{}
	data, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return err
	}
	token, err := jwt.ParseWithClaims(string(data), &uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.Config.JWTSignKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	t.SetClaim(&uc)
	return nil
}

// GetClaim returns token claim
func (t *TokenAuthentication) GetClaim() Claim {
	return t.User.UserClaim
}

// SetClaim sets token claim
func (t *TokenAuthentication) SetClaim(uc Claim) {
	if uc == nil {
		return
	}
	t.User = &UserAuth{
		UserClaim: uc.(interface{}).(*UserClaim),
	}
}
