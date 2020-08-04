package crud

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestUpdateIfSuccessfullyUpdatedName(t *testing.T) {
	s := tests.CreateSuite()

	uid := uint32(100)
	puzzleID := uint(0)
	testName := "updated puzzle name"

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name:   testName,
			UserID: uid,
		},
	}

	expectedData := []models.Puzzle{
		models.Puzzle{
			ID:        1,
			Name:      testName,
			UserID:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	_ = s.Mock.NewRows([]string{"name", "created_at", "updated_at", "user_id"}).
		AddRow(data[0].Name, time.Now(), time.Now(), uid)

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("UPDATE").
		WithArgs(expectedData[0].Name, sqlmock.AnyArg(), puzzleID, expectedData[0].UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	rowsUpdated, err := repo.Update(uid, data[0])

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
