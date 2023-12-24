package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gorm_tutorial/models"
	"gorm_tutorial/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

func GetAllUsers(db *gorm.DB) []models.User {
	var users []models.User
	err := db.Find(&users).Error
	// err := db.Order("username ASC").Find(&users).Error
	// err := db.Order("username ASC, id DESC").Find(&users).Error
	if err != nil {
		log.Fatalln("Cannot get all users:", err)
	}
	return users
}

func GetUserByID(db *gorm.DB, userId uint) (*models.User, error) {
	var user models.User
	err := db.First(&user, userId).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUsersByCreatedAtBefore returns a slice of users with the created_at field less than the given time.
func GetUsersByCreatedAtBefore(db *gorm.DB, t time.Time) ([]models.User, error) {
	var users []models.User
	err := db.Where("created_at < ?", t).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByUsernameAndPassword returns a user with the given username and password.
func GetUserByUsernameAndPassword(db *gorm.DB, username, password string) (models.User, error) {
	var user models.User
	err := db.Where("username = ? AND passwd = md5(?)", username, password).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	return user, nil
}

// InsertUser inserts a new user into the database.
func InsertUser(db *gorm.DB, username, password, email string) (*models.User, error) {
	// Create a new user object.
	user := &models.User{
		Username:      username,
		Password:      fmt.Sprintf("%x", md5.Sum([]byte(password))),
		Email:         email,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Status:        models.UserStatusInactive,
	}

	// Insert the user into the database.
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}

	// Return the user.
	return user, nil
}

// UpdateUserWithChangedPasswordAndEmail updates a user with changed password (MD5) and email.
func UpdateUserWithChangedPasswordAndEmail(db *gorm.DB, userId uint, password, email string) (*models.User, error) {
	// Load the user from the database.
	var user models.User
	err := db.First(&user, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("No user found")
			return nil, nil
		}
		return nil, err
	}

	// Change the password, email, last update.
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user.Email = email
	user.LastUpdatedAt = time.Now()

	// Save the user to the database.
	err = db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	// Return the updated user.
	return &user, nil
}

func GetAllUsersWithPagination(db *gorm.DB, page uint) ([]models.User, error) {
	// Calculate the limit and offset values.
	var limit uint = 5
	offset := (page - 1) * limit

	// Get the total number of users.
	var totalUsersCount int64 = 0
	err := db.Model(&models.User{}).Count(&totalUsersCount).Error
	if err != nil {
		return nil, err
	}

	// Check if the page is valid.
	if page > 0 && page <= uint(totalUsersCount)/limit+1 {
		// Get the users for the specified page.
		var users []models.User
		err := db.Limit(int(limit)).Offset(int(offset)).Find(&users).Error
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	// If the page is invalid, return an empty slice.
	return []models.User{}, nil
}

func GetAllUsersByUsernameLike(db *gorm.DB, inputString string) ([]models.User, error) {
	var users []models.User

	// Construct the SQL query.
	sql := `SELECT * FROM users WHERE username LIKE ?`

	// Execute the query.
	err := db.Raw(sql, "%"+inputString+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetCountUsersGroupedByStatus(db *gorm.DB) ([]models.UserCountGroupedByStatus, error) {
	var userCounts []models.UserCountGroupedByStatus

	// Execute the query.
	err := db.Raw("SELECT user_status, COUNT(*) AS count FROM users GROUP BY user_status").Find(&userCounts).Error
	if err != nil {
		return nil, err
	}

	return userCounts, nil
}

func GetAllUsersWithUsernameLikeUsingProcedure(db *gorm.DB, inputString string) ([]models.User, error) {
	var users []models.User

	// Execute the procedure and scan the results into the users slice.
	err := db.Raw(`CALL get_all_users_with_username_like(?)`, inputString).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllUsersWithTime(db *gorm.DB, inputTime *time.Time) ([]models.User, time.Time, error) {
	// Get the max last_updated_at of all users.
	var lastedTime time.Time
	err := db.Model(&models.User{}).Select("MAX(last_updated_at) AS lasted_time").Scan(&lastedTime).Error
	if err != nil {
		return nil, time.Time{}, err
	}

	// Check if the input time is nil or less than the lasted time.
	if inputTime == nil || inputTime.Before(lastedTime) {
		// Return all users and the lasted time.
		var users []models.User
		err := db.Find(&users).Error
		if err != nil {
			return nil, time.Time{}, err
		}

		return users, lastedTime, nil
	} else {
		return nil, lastedTime, utils.ErrNoModified
	}

}
