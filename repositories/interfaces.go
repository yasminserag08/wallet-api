package repositories

import (
	"wallet-api/models"

	"gorm.io/gorm"
)

// Used for authentication handlers
type UserRepositoryInterface interface {
	CreateUserWithWallet(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

type WalletRepositoryInterface interface {
	Create(userID uint) error
	GetByUserID(userID uint) (models.Wallet, error)
	UpdateBalance(tx *gorm.DB, walletID uint, newBalance int) error
	GetByUserIDForUpdate(tx *gorm.DB, userID uint) (models.Wallet, error)
}

type TransactionRepositoryInterface interface {
	Create(tx *gorm.DB, transaction models.Transaction) error
}
