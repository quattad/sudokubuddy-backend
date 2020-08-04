package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/quattad/sudokubuddy-backend/src/api/utils"
)

// ========= BEFORESAVE() ========= //

// ========= PrepareUser() ========= //
func TestPrepareIfSuccessful(t *testing.T) {

	testUser := User{
		Username:  "johndoe ",
		Email:     " johndoe@gmail.com",
		FirstName: "John ",
		LastName:  "Doe ",
		Password:  " 123456",
	}

	expectedUser := User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Execute function to bested
	testUser.PrepareUser()

	expectedDataBytes, err := json.Marshal(expectedUser)

	if err != nil {
		t.Fatal(err)
	}

	expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	actualDataBytes, err := json.Marshal(testUser)

	if err != nil {
		t.Fatal(err)
	}

	actualDataReader := bytes.NewBuffer(actualDataBytes)

	if err != nil {
		t.Fatal(err)
	}

	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Data does not match. actual: %s, expected: %s", expectedDataBytes, actualDataBytes)
	}

}

// ========= ValidateUser() ========= //
func TestIfValidateUserSuccessfulForDefault(t *testing.T) {
	testUser := User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	// Execute test function
	err := testUser.ValidateUser("")

	if err != nil {
		t.Errorf("Error: %s, expected nil", err)
	}

}

func TestIfValidateUserUserUnsuccessfulForDefaultMissingUsername(t *testing.T) {
	testUser := User{
		Username:  "",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	expectedErr := errors.New("User must have defined property 'username'")

	// Execute test function
	err := testUser.ValidateUser("")

	if err == nil {
		t.Errorf("Error: %s, expected nil", err)
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Actual error: %s, expected error: %s", err, expectedErr)
	}

}

func TestIfValidateUserUnsuccessfulForDefaultEmptyEmail(t *testing.T) {
	testUser := User{
		Username:  "johndoe",
		Email:     "",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	expectedErr := errors.New("User must have defined property 'email'")

	// Execute test function
	err := testUser.ValidateUser("")

	if err == nil {
		t.Errorf("Error: %s, expected nil", err)
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Actual error: %s, expected error: %s", err, expectedErr)
	}

}

func TestIfValidateUserUnsuccessfulForDefaultInvalidEmail(t *testing.T) {
	testUser := User{
		Username:  "johndoe",
		Email:     "invalidemail@.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	expectedErr := errors.New("Invalid email")

	// Execute test function
	err := testUser.ValidateUser("")

	if err == nil {
		t.Errorf("Error: %s, expected nil", err)
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("Actual error: %s, expected error: %s", err, expectedErr)
	}

}

func TestIfValidateUserSuccessfulForLogin(t *testing.T) {
	testUser := User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	// Execute test function
	err := testUser.ValidateUser("login")

	if err != nil {
		t.Errorf("Error: %s, expected nil", err)
	}

}

func TestIfValidateUserSuccessfulForUpdate(t *testing.T) {
	testUser := User{
		Username:  "johndoe",
		Email:     "johndoe@gmail.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "123456",
	}

	// Execute test function
	err := testUser.ValidateUser("update")

	if err != nil {
		t.Errorf("Error: %s, expected nil", err)
	}

}
