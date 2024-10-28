package repositories

import (
	"errors"
	"server/pkg/db"
	"server/pkg/models"
	"time"
)

/*
This auth system only supports one logged in user per account

to THEORETICALLY fix this:
- store refresh tokens in extra database not only in just one field
- query through all access tokens to validate
*/

const (
	RefreshTokenDuration = 7 * 24 * time.Hour
)

func DoesEmailExist(email string) bool {
	var exists bool
	err := db.DB.Model(&models.User{}).Select("1").
		Where("email = ?", email).Limit(1).
		Scan(&exists).Error
	if err != nil {
		return false
	}
	return exists
}

func DoesUserIDExist(userID uint) bool {
	var exists bool
	err := db.DB.Model(&models.User{}).Select("1").
		Where("id = ?", userID).Limit(1).
		Scan(&exists).Error
	if err != nil {
		return false
	}
	return exists
}

func CreateUser(user *models.User) RepositoryResult {
	exists := DoesEmailExist(user.Email)
	if exists {
		return RepositoryResult{
			Result: nil,
			Error:  errors.New("email is already in use"),
		}
	}

	err := db.DB.Save(user).Error
	return RepositoryResult{
		Result: user,
		Error:  err,
	}
}

func GetUserByEmail(email string) UserRepositoryResult {
	var user models.User
	err := db.DB.Model(&models.User{}).Where("email = ?", email).Select("email", "password", "refreshtoken").Find(&user).Error
	if err != nil {
		return UserRepositoryResult{
			Result: &models.User{},
			Error:  err,
		}
	}

	return UserRepositoryResult{
		Result: &user,
		Error:  nil,
	}
}

func GetUserByID(userID uint) UserRepositoryResult {
	var user models.User
	err := db.DB.Model(&models.User{}).Where("id = ?", userID).Select("email", "password", "refreshtoken").Find(&user).Error
	if err != nil {
		return UserRepositoryResult{
			Result: &models.User{},
			Error:  err,
		}
	}

	return UserRepositoryResult{
		Result: &user,
		Error:  nil,
	}
}

func ValidateRefreshToken(userID uint, refreshToken string) bool {
	var user models.User
	err := db.DB.Model(&models.User{}).Where("id = ? AND refreshtoken = ? and refreshtokenexpiry > ?", userID, refreshToken, time.Now()).First(&user).Error
	return err == nil
}

func StoreRefreshToken(userID uint, refreshToken string) error {
	exists := DoesUserIDExist(userID)
	if !exists {
		return errors.New("user does not exist")
	}

	return db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(models.User{RefreshToken: refreshToken, RefreshTokenExpiry: time.Now().Add(RefreshTokenDuration)}).Error
}

func DeleteRefreshToken(userID uint, token string) error {
	err := db.DB.Where("id = ? AND token = ?", userID, token).Updates(models.User{RefreshToken: "", RefreshTokenExpiry: time.Now()}).Error
	return err
}
