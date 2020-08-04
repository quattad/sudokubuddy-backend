package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestUpdateBoardIfSuccessful(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()

	puzzleID := uint32(445)
	expectedValue := 6
	boardRow := 5
	boardCol := 4
	uid := uint32(100)

	// Populate DB and define expected response
	data := []models.Board{
		models.Board{
			Value: expectedValue,
		},
	}

	actualDataBytes, err := json.Marshal(data[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	// 'UPDATE `boards` SET `updated_at` = ?, `value` = ?  WHERE (puzzle_id=? AND board_row=? AND board_col=?)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("UPDATE").
		WithArgs(sqlmock.AnyArg(), expectedValue, puzzleID, boardRow, boardCol).
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
	req, err := http.NewRequest("PUT", "/boards", bytes.NewBuffer(actualDataBytes))

	if err != nil {
		t.Fatal(err)
	}

	// For testing purposes
	// req = mux.SetURLVars(req, map[string]string{
	// 	"puzzle_id": strconv.Itoa(int(puzzleID)),
	// 	"board_row": strconv.Itoa(int(boardRow)),
	// 	"board_col": strconv.Itoa(int(boardCol)),
	// })

	q := req.URL.Query()
	q.Add("puzzle_id", strconv.Itoa(int(puzzleID)))
	q.Add("board_row", strconv.Itoa(int(boardRow)))
	q.Add("board_col", strconv.Itoa(int(boardCol)))
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	UpdateBoard(rr, req)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
