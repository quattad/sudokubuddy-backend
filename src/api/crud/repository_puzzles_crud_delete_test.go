package crud

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestIfDeletePuzzleWasSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	uid := uint32(100)
	puzzleID := uint32(1)
	testName := "updated puzzle name"

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name:      testName,
			UserID:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	_ = s.Mock.NewRows([]string{"name", "user_id"}).
		AddRow(data[0].Name, uid)

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("DELETE").
		WithArgs(puzzleID, data[0].UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	rowsUpdated, err := repo.Delete(puzzleID, uid)

	if err != nil {
		t.Errorf("Error: %s, expected nil", err)
	}

	// Check status code and body
	if rowsUpdated != 1 {
		t.Errorf("Actual rowsUpdated: %v, expected 1", rowsUpdated)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
