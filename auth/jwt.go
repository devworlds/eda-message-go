package auth

import "github.com/devworlds/eda-message-go/auth/internal/jwt"

// ValidateJWT is a public wrapper for the internal jwt.ValidateJWT function.
func ValidateJWT(tokenString string) bool {
	return jwt.ValidateJWT(tokenString)
}
