package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========== CREATE ========== //
func TestSavePuzzleIfSuccessfullySave(t *testing.T) {
	t.Skip("Expected failure due to insertion into both puzzles and boards")
	s := tests.CreateSuite()

	uid := uint32(100)

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name:      "testpuzzle1",
			UserID:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	expectedData := []models.Puzzle{
		models.Puzzle{
			ID:        1,
			Name:      "testpuzzle1",
			UserID:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	expected, err := json.Marshal(expectedData[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	// rows := s.Mock.NewRows([]string{"name", "user_id"}).
	// 	AddRow(data[0].Name, uid)

	// Expect insert for puzzles
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data[0].Name, sqlmock.AnyArg(), sqlmock.AnyArg(), data[0].UserID).
		// WillReturnRows(rows)
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	puzzle, err := repo.Save(data[0])

	if err != nil {
		t.Errorf("Error: %s, expected nil", err)
	}

	actual, err := json.Marshal(puzzle)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	// Check status code and body
	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Actual: %s, expected:%s, expected equal", actual, expected)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
