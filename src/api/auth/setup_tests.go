package auth

import (
	"github.com/jinzhu/gorm"
)

var (
	mockConnect        func(string, string) (*gorm.DB, error)
	mockVerifyPassword func(string, string) error
)

type mockDB struct{}

func (m *mockDB) Connect(DBDRIVER, DBURL string) (*gorm.DB, error) {
	return mockConnect(DBDRIVER, DBURL)
}

type mockSecurity struct{}

func (m *mockSecurity) VerifyPassword(inputPassword string, actualPassword string) error {
	return mockVerifyPassword(inputPassword, actualPassword)
}

func (m *mockSecurity) Hash(string) ([]byte, error) {
	return []byte(""), nil
}
