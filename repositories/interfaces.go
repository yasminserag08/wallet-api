package repositories

import "wallet-api/models"

// Used for authentication handlers
type UserRepositoryInterface interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}
