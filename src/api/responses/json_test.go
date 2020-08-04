package responses

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quattad/sudokubuddy-backend/src/api/utils"
	"github.com/stretchr/testify/require"
)

// ========== JSON() ========== //
func TestJSONIfSuccessfulWrite(t *testing.T) {
	rr := httptest.NewRecorder()

	expectedStatusCode := http.StatusOK

	expectedData := struct {
		Name string
		Num  int
	}{
		Name: "TestName",
		Num:  1,
	}

	expectedDataBytes, err := json.Marshal(expectedData)
	expectedDataString := string(expectedDataBytes)

	if err != nil {
		t.Fatal(err)
	}

	JSON(rr, expectedStatusCode, expectedData)

	// Check responseRecorder
	if actualStatusCode := rr.Code; actualStatusCode != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", actualStatusCode, expectedStatusCode)
	}

	actualDataString := rr.Body.String()
	require.JSONEq(t, actualDataString, expectedDataString)

}

func TestJSONIfFailedWriteWrongStatusCode(t *testing.T) {
	rr := httptest.NewRecorder()

	expectedStatusCode := http.StatusOK

	expectedData := struct {
		Name string
		Num  int
	}{
		Name: "TestName",
		Num:  1,
	}

	JSON(rr, http.StatusUnauthorized, expectedData)

	// Check responseRecorder
	if actualStatusCode := rr.Code; actualStatusCode == expectedStatusCode {
		t.Errorf("Error: handler returned no errors, should return error")
	}
}

func TestJSONIfJSONDoesNotMatch(t *testing.T) {
	rr := httptest.NewRecorder()

	expectedStatusCode := http.StatusOK

	expectedData := struct {
		Name string
		Num  int
	}{
		Name: "TestName",
		Num:  1,
	}

	expectedDataBytes, err := json.Marshal(expectedData)

	if err != nil {
		t.Fatal(err)
	}

	expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	failedData := struct {
		Name string
		Num  int
	}{
		Name: "WrongTestName",
		Num:  10,
	}

	failedDataBytes, err := json.Marshal(failedData)

	if err != nil {
		t.Fatal(err)
	}

	failedDataReader := bytes.NewBuffer(failedDataBytes)

	if err != nil {
		t.Fatal(err)
	}

	JSON(rr, expectedStatusCode, failedData)

	// Check responseRecorder
	if actualStatusCode := rr.Code; actualStatusCode != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", actualStatusCode, expectedStatusCode)
	}

	if utils.JSONEqual(expectedDataReader, failedDataReader) {
		t.Errorf("No error, expected error")
	}

}

// ========== ERROR() ========== //
func TestERRORIfErrorSuccessfullyWritten(t *testing.T) {
	rr := httptest.NewRecorder()
	expectedStatusCode := http.StatusUnauthorized
	expectedError := errors.New("Invalid token")

	expectedData := struct {
		Error string `json:"error"`
	}{
		Error: expectedError.Error(),
	}

	expectedDataBytes, err := json.Marshal(expectedData)

	if err != nil {
		t.Fatal(err)
	}

	expectedDataReader := bytes.NewBuffer(expectedDataBytes)

	ERROR(rr, expectedStatusCode, expectedError)

	actualDataReader := bytes.NewBuffer(rr.Body.Bytes())

	// Check status code and body
	if actualStatusCode := rr.Code; actualStatusCode != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", actualStatusCode, expectedStatusCode)
	}

	if !utils.JSONEqual(expectedDataReader, actualDataReader) {
		t.Errorf("Error: actual and expected io Readers are different, expected same")
	}
}

func TestERRORIfMissingError(t *testing.T) {
	rr := httptest.NewRecorder()
	expectedStatusCode := http.StatusBadRequest

	ERROR(rr, expectedStatusCode, nil)

	// Check status code and body
	if actualStatusCode := rr.Code; actualStatusCode != expectedStatusCode {
		t.Errorf("Error: handler returned status code: %v, expected: %v", actualStatusCode, expectedStatusCode)
	}
}
