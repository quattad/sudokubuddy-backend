package crud

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/quattad/sudokubuddy-backend/src/api/models"
	"github.com/quattad/sudokubuddy-backend/src/api/utils/channels"
)

// UsersCRUDService is a global variable that exposes the methods of the module
var UsersCRUDService UsersCRUDInterface

func init() {
	UsersCRUDService = &UsersCRUD{}
}

// UsersCRUD is a struct that makes it easy to access the db for that particular repository
// by calling r.db
type UsersCRUD struct {
	db *gorm.DB
}

// UsersCRUDInterface is an interface for UsersCRUD struct to allow for mocking of functions
// during testing
type UsersCRUDInterface interface {
	// InitDB
	NewUsersCRUD(*gorm.DB) *UsersCRUD

	// Create
	Save(models.User) (models.User, error)

	// Read
	FindByID(uint32) (models.User, error)
	// FindAll() ([]models.User, error)

	// Update
	// Update(uint32, models.User) (models.User, error)

	// Delete
	// Delete(uint32) (uint32, error)
}

// NewUsersCRUD takes in db as an argument and returns a RepositoryUsersCRUD struct that
// has r.db as a property; making it easy to access the db
func (u *UsersCRUD) NewUsersCRUD(db *gorm.DB) *UsersCRUD {
	u.db = db
	return u
}

// ========== CREATE ========== //

// Save takes a User model and saves it to the db
// Returns the saved model and error if successful, returns empty User instance and error if unsuccessful
func (u *UsersCRUD) Save(user models.User) (models.User, error) {

	var err error
	done := make(chan bool)

	go func(ch chan<- bool) {
		err = u.db.Debug().Model(&models.User{}).Create(&user).Error

		if err != nil {
			ch <- false
		}

		ch <- true
		close(ch)

	}(done)

	if channels.OK(done) {
		return user, nil
	}

	return models.User{}, err

}

// ========== READ ========== //

// FindByID takes a userID and fetches the model instance from the db
// Returns the saved model and error if successful, returns empty User instance and error if unsuccessful
func (u *UsersCRUD) FindByID(uid uint32) (models.User, error) {
	var err error
	user := models.User{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = u.db.Debug().Model(&models.User{}).Where("id=?", uid).Take(&user).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// User found
	if channels.OK(done) {
		return user, nil
	}

	// User not found
	if gorm.IsRecordNotFoundError(err) {
		return user, errors.New("User not found")
	}

	// Other errors
	return user, err

}

// FindAll fetches all the entries from the User model in the db
// Returns an array of models and error if successful, returns empty array and error if unsuccessful
func (u *UsersCRUD) FindAll() ([]models.User, error) {
	var err error
	users := []models.User{}
	done := make(chan bool)

	// If not found or any kind of error, will return false
	go func(ch chan<- bool) {
		defer close(ch)
		err = u.db.Debug().Model(&models.User{}).Limit(100).Find(&users).Error

		if err != nil {
			ch <- false
		}

		ch <- true
	}(done)

	// User found
	if channels.OK(done) {
		return users, nil
	}

	// Other errors
	return nil, err
}

// ========== UPDATE ========== //

// Update takes in a model instance with updated fields and ID, and updates the existing entry in the db that matches ID
// Returns updated model and error if successful, returns empty array and error if unsuccessful
func (u *UsersCRUD) Update(uid uint32, user models.User) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = u.db.Debug().Model(&models.User{}).Where("id=?", uid).UpdateColumns(
			map[string]interface{}{
				"username":   user.Username,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"updated_at": time.Now(),
			},
		)

		ch <- true
	}(done)

	if !channels.OK(done) || rs.Error != nil {
		return 0, rs.Error
	}

	return rs.RowsAffected, nil

}

// ========== DELETE ========== //

// Delete takes in an ID and deletes the existing entry in the db that matches ID
// Returns updated model and error if successful, returns empty array and error if unsuccessful
func (u *UsersCRUD) Delete(uid uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)

		rs = u.db.Debug().Model(&models.User{}).Where("id=?", uid).Delete(&models.User{})
		ch <- true

	}(done)

	if channels.OK(done) && rs.Error == nil {
		return rs.RowsAffected, nil
	}

	return 0, rs.Error

}
