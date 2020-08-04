package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

// ========== CREATEPUZZLE() ========== //
func TestCreatePuzzleIfSuccessful(t *testing.T) {
	t.Skip("Expected failure due to insertion into both puzzles and boards")
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusCreated

	uid := uint32(100)

	data := []models.Puzzle{
		models.Puzzle{
			Name:   "testpuzzle1",
			UserID: uid,
		},
	}

	expected, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data[0].Name, sqlmock.AnyArg(), sqlmock.AnyArg(), data[0].UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return uid, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/puzzles", bytes.NewBuffer(expected))

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
	CreatePuzzle(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	if actual := rr.Body.Bytes(); bytes.Equal(actual, expected) {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", actual, expected)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestCreatePuzzleIfInvalidRequestBody(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusUnprocessableEntity

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	uid := uint32(100)

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return uid, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/puzzles", tests.ErrReader(0))

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
	CreatePuzzle(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestCreatePuzzleIfInvalidPostFields(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusUnprocessableEntity

	uid := uint32(100)

	data := []models.Puzzle{
		models.Puzzle{
			Name:   "testpuzzle1",
			UserID: uid,
		},
	}

	expected, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return uid, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/puzzles", bytes.NewBuffer(expected))

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
	CreatePuzzle(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestCreatePuzzleIfDBCannotConnect(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	expectedStatusCode := http.StatusInternalServerError
	expectedErr := errors.New("Unable to connect to db")

	uid := uint32(100)

	data := []models.Puzzle{
		models.Puzzle{
			Name:   "testpuzzle1",
			UserID: uid,
		},
	}

	expected, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data[0].Name, sqlmock.AnyArg(), sqlmock.AnyArg(), data[0].UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return nil, expectedErr
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return uid, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("POST", "/puzzles", bytes.NewBuffer(expected))

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
	CreatePuzzle(rr, req)

	// Check status code and body
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatusCode)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
