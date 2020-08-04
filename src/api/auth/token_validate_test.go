package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
)

// ========== VALIDATETOKEN() ========== //
func TestValidateTokenIfSuccessful(t *testing.T) {
	config.SECRETKEY = []byte("8KU7Mty9qwQCpKGONftZjkZFB49VoSiG")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = 100
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 60min expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.SECRETKEY)

	req, err := http.NewRequest("GET", "/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Set URL field
	q := req.URL.Query()          // parses raw queries and returns values to a copy of values, not reference
	q.Add("token", tokenString)   // add new value to set
	req.URL.RawQuery = q.Encode() // Encode and reassign to original query

	err = TokenService.ValidateToken(req)

	if err != nil {
		t.Errorf("Error: %v, expected nil", err)
	}

}
