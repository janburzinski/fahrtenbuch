package repositories

import (
	"server/pkg/db"
	"server/pkg/models"
)

func CreateUser(user *models.User) RepositoryResult {
	err := db.DB.Save(user).Error
	return RepositoryResult{
		Result: user,
		Error:  err,
	}
}
