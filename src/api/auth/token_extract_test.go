package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
)

// ========== EXTRACTTOKEN() ========== //
func TestExtractTokenIfSuccessfulByURLQuery(t *testing.T) {
	testUserID := uint32(1001)
	testSecretKey := []byte("DvFb3MYehzqUI8KIDVLvfTU3S1PLwYBL")
	testBearerToken := "Bearer abc123"

	config.SECRETKEY = testSecretKey

	testClaims := jwt.MapClaims{
		"authorized": true,
		"user_id":    testUserID,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	}

	testToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		testClaims,
	)

	testTokenString, err := testToken.SignedString(testSecretKey)

	if err != nil {
		t.Fatal(err)
	}

	// Create request and response objects
	req, err := http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Set URL field
	q := req.URL.Query()            // parse raw queries and returns values to copy of values, not reference
	q.Add("token", testTokenString) // add new value to set
	req.URL.RawQuery = q.Encode()   // encode and reassign to original query

	expectedToken := req.URL.Query().Get("token")

	// Set header
	req.Header.Set("Authorization", testBearerToken)

	actualToken := TokenService.ExtractToken(req)

	if expectedToken != actualToken {
		t.Errorf("Tokens do not match, actual token: %s, expected: %s", expectedToken, actualToken)
	}
}

func TestExtractTokenIfSuccessfulByAuthorizationHeader(t *testing.T) {
	testUserID := uint32(1001)
	testSecretKey := []byte("DvFb3MYehzqUI8KIDVLvfTU3S1PLwYBL")

	config.SECRETKEY = testSecretKey

	testClaims := jwt.MapClaims{
		"authorized": true,
		"user_id":    testUserID,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	}

	testToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		testClaims,
	)

	testTokenString, err := testToken.SignedString(testSecretKey)
	testBearerToken := fmt.Sprintf("Bearer: %s", testTokenString)
	expectedToken := testTokenString

	if err != nil {
		t.Fatal(err)
	}

	// Create request and response objects
	req, err := http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	// Set URL field
	q := req.URL.Query()          // parse raw queries and returns values to copy of values, not reference
	q.Add("token", "")            // add new value to set
	req.URL.RawQuery = q.Encode() // encode and reassign to original query

	// Set header
	req.Header.Set("Authorization", testBearerToken)

	actualToken := TokenService.ExtractToken(req)

	if expectedToken != actualToken {
		t.Errorf("Tokens do not match, actual token: %s, expected: %s", expectedToken, actualToken)
	}
}

// ========== EXTRACTTOKENID() ========== //
func TestExtractTokenIDIfRequestHasValidToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	testUserID := uint32(1001)
	testSecretKey := []byte("DvFb3MYehzqUI8KIDVLvfTU3S1PLwYBL")
	testBearerToken := "Bearer abc123"

	config.SECRETKEY = testSecretKey

	testClaims := jwt.MapClaims{
		"authorized": true,
		"user_id":    testUserID,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	}

	testToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		testClaims,
	)

	testTokenString, err := testToken.SignedString(testSecretKey)

	if err != nil {
		t.Fatal(err)
	}

	// Create request and response objects
	req, err = http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Set URL field
	q := req.URL.Query()            // parse raw queries and returns values to copy of values, not reference
	q.Add("token", testTokenString) // add new value to set
	req.URL.RawQuery = q.Encode()   // encode and reassign to original query

	// Set header
	req.Header.Set("Authorization", testBearerToken)

	extractedID, err := TokenService.ExtractTokenID(req)

	if extractedID != testUserID {
		t.Errorf("ID do not match, ExtractedID: %v, expected: %v", extractedID, testUserID)
	}

}
