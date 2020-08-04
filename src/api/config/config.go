package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// PORT stores port number stored in environment variable
// DBURL stores database URL
// DBDRIVER stores type of db used - in this case MySQL
// SECRETKEY stores the hash key of the API used to generate the jwt
var (
	PORT      int
	DBDRIVER  string
	DBURL     string
	SECRETKEY []byte
)

// Load fetches environment variables and assigns them to respective variables
func Load() {
	// os.ExpandEnv url must be absolute path
	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/github.com/quattad/sudokubuddy-backend/env/.env"))

	if err != nil {
		log.Fatal(err)
	}

	PORT, err = strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		log.Println(err)
		PORT = 9000
	}

	DBDRIVER = os.Getenv("DB_DRIVER")
	DBURL = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	SECRETKEY = []byte(os.Getenv("API_SECRET"))
}
