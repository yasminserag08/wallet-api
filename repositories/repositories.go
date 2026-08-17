package repositories

import (
	"wallet-api/errors"
	"wallet-api/models"

	"gorm.io/gorm"
)

// this repo is used for auth
type UserRepositoryInterface interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

// actual struct
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	user := models.User{}
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return user, errors.ErrNotFound
	}
	return user, result.Error
}

func (r *UserRepository) GetUserByID(id uint) (models.User, error) {
	user := models.User{}
	result := r.db.First(&user, id)
	if result.Error == gorm.ErrRecordNotFound {
		return user, errors.ErrNotFound
	}
	return user, result.Error
}
