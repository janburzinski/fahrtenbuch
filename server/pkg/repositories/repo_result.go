package repositories

import "server/pkg/models"

type RepositoryResult struct {
	Result interface{}
	Error  error
}

type UserRepositoryResult struct {
	Result *models.User
	Error  error
}
