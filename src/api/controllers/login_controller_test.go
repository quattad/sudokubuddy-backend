package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

var (
	mockSignIn   func(string, string) (string, error)
	mockValidate func(string) error
)

// Define mockAuth and methods
type mockAuth struct{}

func (m *mockAuth) SignIn(email string, password string) (string, error) {
	return mockSignIn(email, password)
}

// Define mockUser and methods
type mockUser struct{}

func (u *mockUser) BeforeSave() error {
	return errors.New("placeholder")
}

func (u *mockUser) Prepare() *models.User {
	user := models.User{}
	return &user
}

func (u *mockUser) Validate(action string) error {
	return mockValidate(action)
}

/* =================  Login() ================= */
func TestLoginIfUserLoginSuccessful(t *testing.T) {
	// Should return response.JSON(w *http.ResponseWriter, status 200 OK, token)
	s := tests.CreateSuite()

	testSecretKey := "abcdefgh"
	testEmailOne := "testemailone@gmail.com"
	testPasswordOne := "testpassword123!"
	testExpectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.UK6SDMmQ2lOd0bO1WZg_N7ZIwflOBvGHRkJN-OJYN3Q"

	// Initialize structs with modified interfaces
	auth.AuthService = &mockAuth{}
	config.SECRETKEY = []byte(testSecretKey)

	mockSignIn = func(email, password string) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "1",
		})
		return token.SignedString(config.SECRETKEY)
	}

	// Build recorder for http.ResponseWriter
	rr := httptest.NewRecorder()

	// Build request body for http.Request
	reqBody, err := json.Marshal(map[string]string{
		"email":    testEmailOne,
		"password": testPasswordOne,
	})

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Populate DB
	_ = s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmailOne, testPasswordOne)

	// Execute test function
	LoginControllerService.Login(rr, req)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if actual := rr.Body.String(); actual[1:len(actual)-2] != testExpectedToken {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", actual, testExpectedToken)
	}
}

func TestLoginIfUserLoginUnsuccessful(t *testing.T) {
	// Should return response.ERROR(w, status 401 Unauthorized, err)
	expectedStatusCode := http.StatusUnauthorized
	s := tests.CreateSuite()

	testSecretKey := "abcdefgh"
	testEmailOne := "testemailone@gmail.com"
	testPasswordOne := "testpassword123!"

	// Initialize structs with modified interfaces
	auth.AuthService = &mockAuth{}
	config.SECRETKEY = []byte(testSecretKey)

	mockSignIn = func(email, password string) (string, error) {
		_ = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "1",
		})
		return "", errors.New("Incorrect password")
	}

	rr := httptest.NewRecorder()

	// build request body
	reqBody, err := json.Marshal(map[string]string{
		"email":    testEmailOne,
		"password": testPasswordOne,
	})

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Populate DB with user, email and password
	_ = s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmailOne, testPasswordOne)

	LoginControllerService.Login(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestLoginIfUserLoginInvalidJson(t *testing.T) {
	// Should return response.ERROR(w, status 422 Unprocessable Entity, err)
	s := tests.CreateSuite()

	testSecretKey := "abcdefgh"
	testEmailOne := "testemailone@gmail.com"
	testPasswordOne := "testpassword123!"

	// Initialize structs with modified interfaces
	auth.AuthService = &mockAuth{}
	config.SECRETKEY = []byte(testSecretKey)

	mockSignIn = func(email, password string) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "1",
		})
		return token.SignedString(config.SECRETKEY)
	}

	rr := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "https://golang.org", bytes.NewBuffer([]byte("")))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Populate DB with user, email and password
	_ = s.Mock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, testEmailOne, testPasswordOne)

	LoginControllerService.Login(rr, req)

	// Check status code and body
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusUnprocessableEntity)
	}

}

func TestLoginIfUserLoginInvalidReqBody(t *testing.T) {
	// Should return response.ERROR(w, status 422 Unprocessable Entity, err)
	// s := tests.CreateSuite()
	t.SkipNow()
}

func TestLoginIfUserLoginInvalidUserFields(t *testing.T) {
	// s := tests.CreateSuite()
	t.SkipNow()
}
