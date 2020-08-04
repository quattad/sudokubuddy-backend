package security

import (
	"golang.org/x/crypto/bcrypt"
)

// SecurityService exposes modules for testing
var (
	SecurityService securityServiceInterface
)

func init() {
	SecurityService = &securityService{}
}

type securityService struct {
	Name string
}

type securityServiceInterface interface {
	Hash(string) ([]byte, error)
	VerifyPassword(string, string) error
}

// Hash creates hash from user password and stores it into the database
func (s *securityService) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword uses inbuilt bcrypt method to compare hash stored in db with password given by user
// Return nil on success, error on failure
// Strings are immutable while bytes are immutable, hence need to convert hashedPassword, password into slice of bytes
func (s *securityService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
