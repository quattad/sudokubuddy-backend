package crud

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestUpdateIfUpdateSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	// Populate DB and define expected response
	data := models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	updatedData := models.User{
		Username:  "updateduser",
		Email:     "updatedemail@gmail.com",
		FirstName: "Updatedfirst",
		LastName:  "Updatedsecond",
		Password:  "updatedpassword123!",
	}

	_ = s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data.Username, data.Email, data.FirstName, data.LastName, data.Password)

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := UsersCRUDService.NewUsersCRUD(s.DB)
	rowsUpdated, err := repo.Update(1, updatedData)

	if err != nil {
		t.Errorf("Actual err: %s, expected nil", err)
	}

	if rowsUpdated != 1 {
		t.Errorf("Actual rowsUpdated: %v, expected 1", rowsUpdated)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
