package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/utils/console"
)

// TokenService is variable of tokenServiceInterface type that exposes the
// functions of tokenService to other modules
var TokenService tokenServiceInterface

func init() {
	TokenService = &tokenService{}
}

// tokenServiceInterface has functions CreateToken, ValidateToken, ExtractToken and
// ExtractTokenID
type tokenServiceInterface interface {
	CreateToken(uint32) (string, error)
	ValidateToken(*http.Request) error
	ExtractToken(*http.Request) string
	ExtractTokenID(*http.Request) (uint32, error)
}

type tokenService struct{}

// CreateToken creates a jwt token
func (t *tokenService) CreateToken(userID uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 60min expiry
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.SECRETKEY)
}

// ValidateToken takes in a http.Request object and validates if a token it holds
// and returns nil if no error
func (t *tokenService) ValidateToken(req *http.Request) error {
	/*
		1. Extract token using ExtractToken function
		2. Check if signing method is HS256, if err return err with unexpected signing method
		3. Parse token with config.SECRETKEY, if err return the err
		4. If claims are ok and the token is valid, return nil
	*/
	tokenString := TokenService.ExtractToken(req)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.SECRETKEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		console.Pretty(claims)
		return nil
	}

	return err
}

// ExtractToken takes in a http.Request object and validates if a token it holds
// and returns nil if no error
func (t *tokenService) ExtractToken(req *http.Request) string {
	/*
		1. Extract token from URL query
		2. If token does not exist in URL query, get bearer token from 'Authorization' header
		3a. Bearer token should be of the form 'Bearer <bearerToken>'.
		3b. Parse by splitting by space, and return second element of split.
	*/

	token := req.URL.Query().Get("token")

	if token != "" {
		return token
	}

	bearerTokenString := req.Header.Get("Authorization")
	split := strings.Split(bearerTokenString, " ")

	if len(split) == 1 {
		return ""
	}

	return split[1]

}

// ExtractTokenID takes in a http.Request object and validates if a token it holds
// and returns nil if no error
func (t *tokenService) ExtractTokenID(req *http.Request) (uint32, error) {
	tokenString := TokenService.ExtractToken(req)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.SECRETKEY, nil
	})

	if err != nil {
		return 0, err
	}

	// decoding token.Claims to jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)

		if err != nil {
			return 0, err
		}

		return uint32(uid), nil
	}

	return 0, errors.New("Token not valid")

}
