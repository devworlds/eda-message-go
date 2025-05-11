package auth

// JWTValidatorAdapter adapts the ValidateJWT function to the Validator interface.
type JWTValidatorAdapter struct{}

// ValidateJWT validates a JWT token.
func (v JWTValidatorAdapter) ValidateJWT(token string) bool {
	return ValidateJWT(token)
}
