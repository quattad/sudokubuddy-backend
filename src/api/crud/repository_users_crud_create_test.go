package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/security"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

func TestSaveIfSaveSuccessful(t *testing.T) {
	t.Fatal("Resolve checking hash match")

	s := tests.CreateSuite()
	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()

	// Populate DB and define expected response
	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "123456",
			CreatedAt: expectedCreatedAt,
			UpdatedAt: expectedUpdatedAt,
		},
	}

	generatedHash, err := security.SecurityService.Hash(data[0].Password)

	if err != nil {
		t.Errorf("Error: %s while hashing test password.", generatedHash)
	}

	generatedHashString := string(generatedHash)

	expectedData := []models.User{
		models.User{
			ID:        1,
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  generatedHashString,
			CreatedAt: expectedCreatedAt,
			UpdatedAt: expectedUpdatedAt,
		},
	}

	expected, err := json.Marshal(expectedData[0])

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	expectedDataReader := bytes.NewBuffer(expected)

	// Anyarg() for password since already hashed
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("INSERT INTO").
		WithArgs(data[0].Username, data[0].Email, data[0].FirstName, data[0].LastName, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := UsersCRUDService.NewUsersCRUD(s.DB)
	post, err := repo.Save(data[0])

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(post)

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
