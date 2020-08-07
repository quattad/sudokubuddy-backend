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
	s := tests.CreateSuite()

	expectedCreatedAt := time.Now()
	expectedUpdatedAt := time.Now()
	testPassword := "123456"

	// Populate DB and define expected response
	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  testPassword,
			CreatedAt: expectedCreatedAt,
			UpdatedAt: expectedUpdatedAt,
		},
	}

	// generatedHash, err := security.SecurityService.Hash(testPassword)

	// if err != nil {
	// 	t.Errorf("Error: %s while hashing test password.", generatedHash)
	// }

	// generatedHashString := string(generatedHash)

	// Generate data for everything but password
	expectedData := []models.User{
		models.User{
			ID:        1,
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
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
	user, err := repo.Save(data[0])

	// Create a new user with everything but password
	extractedUser := models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err != nil {
		t.Fatal(err)
	}

	actual, err := json.Marshal(extractedUser)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	actualDataReader := bytes.NewBuffer(actual)

	// Check status code and body
	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Actual: %s, expected:%s, expected equal", actual, expected)
	}

	// Compare passwords
	if err := security.SecurityService.VerifyPassword(user.Password, testPassword); err != nil {
		t.Errorf("Passwords are not equal")
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
