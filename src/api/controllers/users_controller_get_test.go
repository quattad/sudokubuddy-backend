package controllers

// ========== GETUSER() ========== //
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/tests"
	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

/* =================  GETUSER() ================= */
func TestGetUserIfSuccessfulGet(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()
	uid := 1

	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "123456",
		},
	}

	expected, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	rows := s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data[0].Username, data[0].Email, data[0].FirstName, data[0].LastName, data[0].Password)

	s.Mock.ExpectQuery("SELECT *").WithArgs(uid).WillReturnRows(rows)

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("GET", "/users", nil)

	req = mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(uid),
	})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	GetUser(rr, req)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	if actual := rr.Body.Bytes(); bytes.Equal(actual, expected) {
		t.Errorf("Error: handler returned unexpected body: %v, expected: %v", actual, expected)
	}
}

// ========== GETUSERS() ========== //
func TestGetUsersIfSuccessfulGet(t *testing.T) {
	// Populate DB
	s := tests.CreateSuite()

	data := []models.User{
		models.User{
			Username:  "johndoe",
			Email:     "johndoe@gmail.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "123456",
		},
		models.User{
			Username:  "timdoe",
			Email:     "timdoe@gmail.com",
			FirstName: "Tim",
			LastName:  "Doe",
			Password:  "09876432",
		},
	}

	expected, err := json.Marshal(data)
	expectedDataReader := bytes.NewBuffer(expected)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	rows := s.Mock.NewRows([]string{"username", "email", "first_name", "last_name", "password"}).
		AddRow(data[0].Username, data[0].Email, data[0].FirstName, data[0].LastName, data[0].Password).
		AddRow(data[1].Username, data[1].Email, data[1].FirstName, data[1].LastName, data[1].Password)

	s.Mock.ExpectQuery("SELECT *").WillReturnRows(rows)

	// Initialize struct with modified interfaces
	database.DBService = &dbMock{}

	// Define custom functions
	mockConnect = func(DBDRIVER, DBURL string) (*gorm.DB, error) {
		return s.DB, nil
	}

	// Build request and response objects
	req, err := http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	if err != nil {
		t.Fatal(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute function to be tested
	GetUsers(rr, req)

	actualData := rr.Body.Bytes()
	actualDataReader := bytes.NewBuffer(actualData)

	// Check status code and body
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error: handler returned status code: %v, expected: %v", status, http.StatusOK)
	}

	// if !bytes.Equal(actualData, expected) {
	// 	t.Errorf("Error: data does not match, actual body: %s, expected: %s", actualData, expected)
	// }

	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Error: data does not match, actual body: %s, expected: %s", actualData, expected)
	}
}
