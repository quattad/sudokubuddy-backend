package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/auth"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
)

// ========== DELETEUSER() ========== //
func TestDeleteUserIfSuccessful(t *testing.T) {
	s := tests.CreateSuite()

	expectedStatus := http.StatusOK

	// Populate DB
	data := models.User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	_ = s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data.Username, data.Email, data.FirstName, data.LastName, data.Password)

	// Set SQL expectations
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec("DELETE").
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}
	auth.TokenService = &tokenMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	mockExtractTokenID = func(r *http.Request) (uint32, error) {
		return 1, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("DELETE", "/posts", nil)

	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	DeleteUser(rr, req)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, expectedStatus)
	}
}
