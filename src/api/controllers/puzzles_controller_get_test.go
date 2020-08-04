package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========== GETPUZZLE() ========== //
func TestGetPuzzleIfSuccessfulGet(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	uid := uint32(100)
	puzzleID := uint32(125)
	createdAtExpected := time.Now()
	updatedAtExpected := time.Now()

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			ID:        puzzleID,
			Name:      "testpuzzle1",
			UserID:    uid,
			CreatedAt: createdAtExpected,
			UpdatedAt: updatedAtExpected,
		},
	}

	expectedDataBytes, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	rows := s.Mock.NewRows([]string{"id", "name", "created_at", "updated_at", "user_id"}).
		AddRow(puzzleID, data[0].Name, createdAtExpected, updatedAtExpected, uid)

	s.Mock.ExpectQuery("SELECT *").
		WithArgs(puzzleID, uid).
		WillReturnRows(rows)

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
	req, err := http.NewRequest("GET", "/puzzles", nil)

	req = mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(int(puzzleID)),
	})

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
	GetPuzzle(rr, req)

	actualDataBytes := rr.Body.Bytes()
	actualDataReader := bytes.NewBuffer(actualDataBytes)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if !utils.JSONEqual(actualDataReader, expectedDataReader) {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", string(actualDataBytes), string(expectedDataBytes))
	}
}

// ========== GETPUZZLES() ========== //
func TestGetPuzzlesIfSuccessfulGet(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()

	uidOne := uint32(100)
	// uidTwo := uint32(9495)

	puzzleIDOne := uint32(125)
	puzzleIDTwo := uint32(150429)

	createdAtExpected := time.Now()
	updatedAtExpected := time.Now()

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			ID:        puzzleIDOne,
			Name:      "testpuzzle1",
			UserID:    uidOne,
			CreatedAt: createdAtExpected,
			UpdatedAt: updatedAtExpected,
		},
		models.Puzzle{
			ID:        puzzleIDTwo,
			Name:      "testpuzzle2",
			UserID:    uidOne,
			CreatedAt: createdAtExpected,
			UpdatedAt: updatedAtExpected,
		},
	}

	expectedDataBytes, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	rows := s.Mock.NewRows([]string{"id", "name", "created_at", "updated_at", "user_id"}).
		AddRow(puzzleIDOne, data[0].Name, createdAtExpected, updatedAtExpected, uidOne).
		AddRow(puzzleIDTwo, data[1].Name, createdAtExpected, updatedAtExpected, uidOne)

	s.Mock.ExpectQuery("SELECT *").
		WithArgs(uidOne).
		WillReturnRows(rows)

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return uidOne, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("GET", "/puzzles", nil)

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
	GetPuzzles(rr, req)

	actualDataBytes := rr.Body.Bytes()
	actualDataReader := bytes.NewBuffer(actualDataBytes)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if !utils.JSONEqual(actualDataReader, expectedDataReader) {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", string(actualDataBytes), string(expectedDataBytes))
	}
}
