package handler

import "github.com/golang-jwt/jwt/v5"

// AccessTokenClaims is a struct thats used for JWT claims
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email                       string                 `json:"email"`
	Phone                       string                 `json:"phone"`
	AppMetaData                 map[string]interface{} `json:"app_metadata"`
	UserMetaData                map[string]interface{} `json:"user_metadata"`
	Role                        string                 `json:"role"`
	AuthenticatorAssuranceLevel string                 `json:"aal,omitempty"`
	// AuthenticationMethodReference []models.AMREntry      `json:"amr,omitempty"`
	SessionId   string `json:"session_id,omitempty"`
	IsAnonymous bool   `json:"is_anonymous"`
}
