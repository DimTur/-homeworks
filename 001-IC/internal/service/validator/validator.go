package validator

// TokenValidator interface for token validation
//
//go:generate mockgen -source=validator.go -destination=validator_mock.go -package validator
type TokenValidator interface {
	Validate(token string) bool
}

// MainTokenValidator validates tokens against a set of allowed tokens
type MainTokenValidator struct {
	allowedTokens map[string]bool
}

// NewMainTokenValidator creates a new instance of MainTokenValidator
func NewMainTokenValidator(tokens []string) *MainTokenValidator {
	allowedTokens := make(map[string]bool)
	for _, token := range tokens {
		allowedTokens[token] = true
	}
	return &MainTokenValidator{allowedTokens: allowedTokens}
}

// Validate checks if a token is valid
func (mtv *MainTokenValidator) Validate(token string) bool {
	return mtv.allowedTokens[token]
}
