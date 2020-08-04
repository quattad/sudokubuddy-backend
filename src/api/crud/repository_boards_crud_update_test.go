package crud

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestUpdateIfSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	puzzleID := uint32(445)
	expectedValue := 6
	boardRow := 5
	boardCol := 4

	// Populate DB and define expected response
	data := []models.Board{
		models.Board{
			Value:    expectedValue,
			BoardRow: boardRow,
			BoardCol: boardCol,
		},
	}

	// 'UPDATE `boards` SET `updated_at` = ?, `value` = ?  WHERE (puzzle_id=? AND board_row=? AND board_col=?)
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("UPDATE").
		WithArgs(sqlmock.AnyArg(), expectedValue, puzzleID, boardRow, boardCol).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	repo := BoardsCRUDService.NewBoardsCRUD(s.DB)
	rowsUpdated, err := repo.Update(puzzleID, data[0])

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
