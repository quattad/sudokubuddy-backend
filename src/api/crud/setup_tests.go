package crud

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

/* =================  MOCK STRUCTS ================= */
type dbMock struct{}
type tokenMock struct{}

/* =================  MOCK FUNCTION DECLARATIONS ================= */
var (
	mockConnect        func(string, string) (*gorm.DB, error)
	mockCreateToken    func(uint32) (string, error)
	mockValidateToken  func(*http.Request) error
	mockExtractToken   func(*http.Request) string
	mockExtractTokenID func(*http.Request) (uint32, error)
)

/* =================  MOCK CALLERS ================= */
// DBMOCK
func (d *dbMock) Connect(DBDRIVER, DBURL string) (*gorm.DB, error) {
	return mockConnect(DBDRIVER, DBURL)
}

// TOKENMOCK
func (t *tokenMock) CreateToken(uid uint32) (string, error) {
	return mockCreateToken(uid)
}

func (t *tokenMock) ValidateToken(r *http.Request) error {
	return mockValidateToken(r)
}

func (t *tokenMock) ExtractToken(r *http.Request) string {
	return mockExtractToken(r)
}

func (t *tokenMock) ExtractTokenID(r *http.Request) (uint32, error) {
	return mockExtractTokenID(r)
}
