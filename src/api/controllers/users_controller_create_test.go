package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

// ========== CREATEUSER() ========== //
func TestCreateUserIfSuccessful(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusCreated

	data := models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	// Anyarg() for password since already hashed
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data.Username, data.Email, data.FirstName, data.LastName, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(expected))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	CreateUser(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	if actual := rr.Body.Bytes(); bytes.Equal(actual, expected) {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", actual, expected)
	}
}

func TestCreateUserIfInvalidRequestBody(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusUnprocessableEntity

	// data := models.User{
	// 	Username:  "johndoe",
	// 	Email:     "johndoe@gmail.com",
	// 	FirstName: "John",
	// 	LastName:  "Doe",
	// 	Password:  "123456",
	// }

	// expected, err := json.Marshal(data)

	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	// }

	// s.Mock.ExpectBegin()
	// s.Mock.ExpectExec("INSERT INTO").
	// 	WithArgs(data.Username, data.Email, data.FirstName, data.LastName, data.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/users", tests.ErrReader(0))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	CreateUser(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}
}

func TestCreateUserIfInvalidPostFields(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusUnprocessableEntity

	data := models.User{}

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	// s.Mock.ExpectBegin()
	// s.Mock.ExpectExec("INSERT INTO").
	// 	WithArgs(data.Username, data.Email, data.FirstName, data.LastName, data.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(expected))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	CreateUser(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}
}

func TestCreateUserIfDBCannotConnect(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusInternalServerError
	expectedErr := errors.New("Connection to db failed")

	data := models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data.Username, data.Email, data.FirstName, data.LastName, data.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, expectedErr
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(expected))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	CreateUser(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}
}
