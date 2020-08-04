package crud

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

func TestDeleteIfSuccessfulDelete(t *testing.T) {
	s := tests.CreateSuite()
	var uid uint32 = 10

	// Populate DB and define expected response
	data := models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	_ = s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data.Username, data.Email, data.FirstName, data.LastName, data.Password)

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("DELETE").
		WillReturnResult(sqlmock.NewResult(int64(uid), 1))
	s.Mock.ExpectCommit()

	// Execute function to be tested
	// Cannot use UsersCRUDService here as it is a UserCRUDInterface instance which does not have 'db' field
	repo := UsersCRUDService.NewUsersCRUD(s.DB)
	rowsDeleted, err := repo.Delete(uid)

	if err != nil {
		t.Errorf("Actual err: %s, expected nil", err)
	}

	if rowsDeleted != 1 {
		t.Errorf("Actual rowsDeleted: %v, expected 1", rowsDeleted)
	}

	// ensure all expectations have been met
	if err = s.Mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}
