package crud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========== FindByID() ========== //
func TestFindByIDIfSuccessfulPuzzle(t *testing.T) {

	s := tests.CreateSuite()

	uid := uint32(1)
	puzzleID := uint32(10)

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			ID:   puzzleID,
			Name: "testpuzzle1",
		},
	}

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"name"}).
		AddRow(data[0].Name)

	s.Mock.ExpectQuery("SELECT *").WithArgs(puzzleID, uid).WillReturnRows(rows)

	// Execute function to be tested
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	puzzle, err := repo.FindByID(puzzleID, uid)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(puzzle)

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

func TestFindByIDIfPuzzleDoesNotExist(t *testing.T) {
	var err error
	s := tests.CreateSuite()
	expectedErr := errors.New("Puzzle not found")

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name: "testpuzzle1",
		},
	}

	uid := uint32(100)
	puzzleID := uint32(10)

	_ = s.Mock.NewRows([]string{"name"}).
		AddRow(data[0].Name)

	s.Mock.ExpectQuery("SELECT *").WithArgs(puzzleID, uid).WillReturnError(expectedErr)

	// Execute function to be tested
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	_, actualErr := repo.FindByID(puzzleID, uid)

	if actualErr == nil {
		t.Errorf("No error, expected error: Puzzle not found")
	}

	if actualErr != expectedErr {
		t.Errorf("Actual err: %s, expected %s", actualErr, expectedErr)
	}

	// ensure all expectations have been met
	if unmetErr := s.Mock.ExpectationsWereMet(); unmetErr != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

// ========== FindAll() ========== //
func TestFindAllIfSuccessful(t *testing.T) {

	s := tests.CreateSuite()

	// Populate DB and define expected response
	data := []models.Puzzle{
		models.Puzzle{
			Name: "testPuzzleOne",
		},
		models.Puzzle{
			Name: "testPuzzleTwo",
		},
	}

	uid := uint32(1001)

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"name"}).
		AddRow(data[0].Name).
		AddRow(data[1].Name)

	s.Mock.ExpectQuery("SELECT *").WithArgs(uid).WillReturnRows(rows)

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := PuzzlesCRUDService.NewPuzzlesCRUD(s.DB)
	puzzles, err := repo.FindAll(uid)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(puzzles)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Actual data and expected data are not the same")
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
