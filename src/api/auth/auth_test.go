package auth

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/security"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

/* ================= SignIn ================= */
func TestSignInIfCorrectEmailAndPassword(t *testing.T) {
	s := tests.CreateSuite()

	testEmail := "testuser@gmail.com"
	testPassword := "testpassword"

	// define mock functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockVerifyPassword = func(inputPassword, actualPassword string) error {
		return nil
	}

	// Initialize structs with modified interfaces
	database.DBService = &mockDB{}
	security.SecurityService = &mockSecurity{}

	// expect required db actions
	rows := s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmail, testPassword)

	// declare expected calls from function
	// if not will throw 'all expectations already fulfilled' err
	const sqlSelectOne = `SELECT *`
	s.Mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(testEmail).
		WillReturnRows(rows)

	_, err := AuthService.SignIn(testEmail, testPassword)

	if err != nil {
		t.Errorf("Error: %v, expected nil", err)
	}
}

func TestSignInIfDatabaseConnectionFailure(t *testing.T) {
	s := tests.CreateSuite()

	testEmail := "testuser@gmail.com"
	testPassword := "testpassword"
	testError := errors.New("Mock db error")

	// define mock functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		// fmt.Println("Execute mockConnect, returning no db with error")
		return s.DB, testError
	}

	mockVerifyPassword = func(inputPassword, actualPassword string) error {
		return nil
	}

	// Initialize structs with modified interfaces
	database.DBService = &mockDB{}
	security.SecurityService = &mockSecurity{}

	_, err := AuthService.SignIn(testEmail, testPassword)

	if err == nil {
		t.Errorf("No error, expected error:%v", testError)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestSignInIfIncorrectEmail(t *testing.T) {
	s := tests.CreateSuite()

	testEmail := "testuser@gmail.com"
	testWrongEmail := "wrongemail@gmail.com"
	testPassword := "testpassword"
	testError := errors.New("Email not found")

	// define mock functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		// fmt.Println("Execute mockConnect, returning mock db with no error")
		return s.DB, nil
	}

	mockVerifyPassword = func(inputPassword, actualPassword string) error {
		// fmt.Println("Execute mockVerifyPassword, returning testError")
		return testError
	}

	// Initialize structs with modified interfaces
	database.DBService = &mockDB{}
	security.SecurityService = &mockSecurity{}

	// expect required db actions
	_ = s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmail, testPassword)

	// declare expected calls from function
	// if not will throw 'all expectations already fulfilled' err
	const sqlSelectOne = `SELECT *`
	s.Mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(testWrongEmail).
		WillReturnError(testError)

	_, err := AuthService.SignIn(testWrongEmail, testPassword)

	if err == nil {
		t.Errorf("No error, expected error: %v", testError)
	}
}

func TestSignInIfIncorrectPassword(t *testing.T) {
	s := tests.CreateSuite()

	testEmail := "testuser@gmail.com"
	testPassword := "testpassword"
	testError := errors.New("Password does not match")

	// define mock functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockVerifyPassword = func(inputPassword, actualPassword string) error {
		return testError
	}

	// Initialize structs with modified interfaces
	database.DBService = &mockDB{}
	security.SecurityService = &mockSecurity{}

	// expect required db actions
	_ = s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmail, testPassword)

	// declare expected calls from function
	// if not will throw 'all expectations already fulfilled' err
	const sqlSelectOne = `SELECT *`
	s.Mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(testEmail).
		WillReturnError(testError)

	_, err := AuthService.SignIn(testEmail, "12345")

	if err == nil {
		t.Errorf("No error, expected error:%v", testError)
	}
}
