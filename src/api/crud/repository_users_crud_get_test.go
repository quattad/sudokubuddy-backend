package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========== FindByID() ========== //
func TestFindByIDIfSuccessful(t *testing.T) {

	s := tests.CreateSuite()

	// Populate DB and define expected response
	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "123456",
		},
	}

	uid := uint32(1)

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	rows := s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data[0].Username, data[0].Email, data[0].FirstName, data[0].LastName, data[0].Password)

	s.Mock.ExpectQuery("SELECT *").WithArgs(uid).WillReturnRows(rows)

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := UsersCRUDService.NewUsersCRUD(s.DB)
	post, err := repo.FindByID(uid)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(post)

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

func TestFindByIDIfUserDoesNotExist(t *testing.T) {
	var err error
	s := tests.CreateSuite()

	// Populate DB and define expected response
	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "123456",
		},
	}

	expectedErr := gorm.ErrRecordNotFound

	uid := uint32(100)

	_ = s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data[0].Username, data[0].Email, data[0].FirstName, data[0].LastName, data[0].Password)

	s.Mock.ExpectQuery("SELECT *").WithArgs(uid).WillReturnError(expectedErr)

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := UsersCRUDService.NewUsersCRUD(s.DB)
	_, actualErr := repo.FindByID(uid)

	if actualErr == nil {
		t.Errorf("No error, expected error: User not found")
	}

	if actualErr == expectedErr {
		t.Errorf("Actual err: %s, expected error: %s", actualErr, expectedErr)
	}

	// ensure all expectations have been met
	if unmetErr := s.Mock.ExpectationsWereMet(); unmetErr != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
