package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
)

// ========== CREATETOKEN() ========== //
func TestCreateTokenIfSuccessful(t *testing.T) {
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

	expectedTokenString, err := testToken.SignedString(testSecretKey)

	if err != nil {
		t.Fatal(err)
	}

	actualTokenString, err := TokenService.CreateToken(testUserID)

	if err != nil {
		t.Fatal(err)
	}

	if expectedTokenString != actualTokenString {
		t.Errorf("Error: actualTokenString: %s, expected: %s", actualTokenString, expectedTokenString)
	}

}
