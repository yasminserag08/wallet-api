package repositories

import (
	"errors"
	appErrors "wallet-api/errors"
	"wallet-api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUserWithWallet(user models.User) (models.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		wallet := models.Wallet{UserID: user.ID, Balance: 0}
		if err := tx.Create(&wallet).Error; err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	user := models.User{}
	result := r.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, appErrors.ErrNotFound
	}
	return user, result.Error
}

func (r *UserRepository) GetUserByID(id uint) (models.User, error) {
	user := models.User{}
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, appErrors.ErrNotFound
	}
	return user, result.Error
}
