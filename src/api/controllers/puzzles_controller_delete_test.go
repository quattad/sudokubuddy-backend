package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

// ========== DELETEPUZZLE() ========== //
func TestDeletePuzzleIfSuccessfulyDeleted(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()

	uid := uint32(100)
	puzzleID := uint32(125)

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name:   "new puzzle name",
			UserID: uid,
		},
	}

	actualDataBytes, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	_ = s.Mock.NewRows([]string{"name", "created_at", "updated_at", "user_id"}).
		AddRow(data[0].Name, time.Now(), time.Now(), data[0].UserID)

	// Set SQL expectations
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("DELETE").
		WithArgs(puzzleID, uid).
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
	req, err := http.NewRequest("DELETE", "/puzzles", bytes.NewBuffer(actualDataBytes))

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
	DeletePuzzle(rr, req)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if actual := rr.Body.Bytes(); bytes.Equal(actual, []byte("1")) {
		t.Errorf("Body was not as expected")
	}
}
