package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========== FINDBYID() ========== //
func TestBoardsFindByIDIfSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	data := []models.Board{
		models.Board{
			BoardRow: 1,
			BoardCol: 1,
			Value:    2,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 2,
			Value:    5,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 3,
			Value:    7,
		},
	}

	expectedData := []models.Board{
		models.Board{
			ID:       2,
			BoardRow: 1,
			BoardCol: 2,
			Value:    5,
		},
	}

	expected, err := json.Marshal(expectedData[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"board_row", "board_col", "value"}).
		AddRow(data[1].BoardRow, data[1].BoardCol, data[1].Value)

	s.Mock.ExpectQuery("SELECT *").
		WithArgs(expectedData[0].ID).
		WillReturnRows(rows)

	// Execute function to be tested
	repo := BoardsCRUDService.NewBoardsCRUD(s.DB)
	board, err := repo.FindByID(expectedData[0].ID)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(board)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	// Check status code and body
	if utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("No error, expected error")
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}

}

// ========== FINDBYPUZZLEIDROWCOL() ========== //
func TestBoardsFindByPuzzleIDRowColIfSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	data := []models.Board{
		models.Board{
			BoardRow: 1,
			BoardCol: 1,
			Value:    2,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 2,
			Value:    5,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 3,
			Value:    7,
		},
	}

	expectedData := []models.Board{
		models.Board{
			ID:       3,
			BoardRow: 1,
			BoardCol: 3,
			Value:    7,
		},
	}

	expected, err := json.Marshal(expectedData[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"board_row", "board_col", "value"}).
		AddRow(data[1].BoardRow, data[1].BoardCol, data[1].Value)

	s.Mock.ExpectQuery("SELECT *").
		WithArgs(expectedData[0].PuzzleID, expectedData[0].BoardRow, expectedData[0].BoardCol).
		WillReturnRows(rows)

	// Execute function to be tested
	repo := BoardsCRUDService.NewBoardsCRUD(s.DB)
	board, err := repo.FindByPuzzleIDRowCol(expectedData[0].PuzzleID, expectedData[0].BoardRow, expectedData[0].BoardCol)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(board)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	// Check status code and body
	if utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("No error, expected error")
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

// ========== FINDALL() ========== //
func TestBoardsFindAllIfSuccessful(t *testing.T) {
	s := tests.CreateSuite()
	testUID := uint32(2)

	data := []models.Board{
		models.Board{
			BoardRow: 1,
			BoardCol: 1,
			Value:    2,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 2,
			Value:    5,
		},
		models.Board{
			BoardRow: 1,
			BoardCol: 3,
			Value:    7,
		},
	}

	expectedData := []models.Board{
		models.Board{
			ID:       2,
			BoardRow: 1,
			BoardCol: 2,
			Value:    5,
		},
	}

	expected, err := json.Marshal(expectedData[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"board_row", "board_col", "value"}).
		AddRow(data[1].BoardRow, data[1].BoardCol, data[1].Value)

	s.Mock.ExpectQuery("SELECT *").
		WithArgs(testUID).
		WillReturnRows(rows)

	// Execute function to be tested
	repo := BoardsCRUDService.NewBoardsCRUD(s.DB)
	board, err := repo.FindAll(testUID)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(board)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	// Check status code and body
	if utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("No error, expected error")
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
