package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/config"
	"github.com/quattad/sudokubuddy-backend/src/api/database"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/security"
	"github.com/quattad/sudokubuddy-backend/src/api/utils/channels"
)

// AuthService is a global variable that exposes the SignIn methods to
// other modules
var (
	AuthService authServiceInterface
)

func init() {
	AuthService = &authService{}
}

type authService struct {
	Name string
}

type authServiceInterface interface {
	SignIn(string, string) (string, error)
}

// SignIn connects to the database, verifies if user with specific email is in the db
// Then checks if the stored hashed password and password provided match
// If matches, then generates a jwt token with user.ID

func (a *authService) SignIn(email, password string) (string, error) {
	var db *gorm.DB
	var err error

	done := make(chan bool)
	user := models.User{}

	go func(ch chan<- bool) {
		db, err = database.DBService.Connect(config.DBDRIVER, config.DBURL)

		if err != nil {
			ch <- false
			return
		}

		err = db.Debug().Model(&models.User{}).Where("email=?", email).Take(&user).Error

		if err != nil {
			ch <- false
			return
		}

		err = security.SecurityService.VerifyPassword(user.Password, password)

		if err != nil {
			ch <- false
			return
		}

		ch <- true
		close(ch)
	}(done)

	if channels.OK(done) {
		return TokenService.CreateToken(user.ID)
	}

	return "", err
}
