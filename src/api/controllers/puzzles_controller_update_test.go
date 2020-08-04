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

// ========== UPDATEPUZZLE() ========== //
func TestUpdatePuzzleIfSuccessfullyUpdateName(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()

	uid := uint32(100)
	puzzleID := uint32(125)

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name:   "updatedpuzzlename",
			UserID: uid,
		},
	}

	actualDataBytes, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	// expectedData := []models.Puzzle{
	// 	models.Puzzle{
	// 		ID:        1,
	// 		Name:      "updatedpuzzlename",
	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 		UserID:    uid,
	// 	},
	// }

	// expectedDataBytes, err := json.Marshal(expectedData[0])

	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	// }

	// expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	_ = s.Mock.NewRows([]string{"name", "created_at", "updated_at", "user_id"}).
		AddRow(data[0].Name, time.Now(), time.Now(), data[0].UserID)

	// Set SQL expectations
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("UPDATE").
		WithArgs(data[0].Name, sqlmock.AnyArg(), puzzleID, data[0].UserID).
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
	req, err := http.NewRequest("PUT", "/puzzles", bytes.NewBuffer(actualDataBytes))

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
	UpdatePuzzle(rr, req)

	// actualDataBytes := rr.Body.Bytes()
	// actualDataReader := bytes.NewBuffer(actualDataBytes)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if actual := rr.Body.Bytes(); bytes.Equal(actual, []byte("1")) {
		t.Errorf("Body was not as expected")
	}
}
