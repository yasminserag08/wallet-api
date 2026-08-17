package repositories

import (
	"errors"
	appErrors "wallet-api/errors"
	"wallet-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (w *WalletRepository) Create(userID uint) error {
	result := w.db.Create(&models.Wallet{
		UserID:  userID,
		Balance: 0,
	})

	return result.Error
}

func (w *WalletRepository) GetByUserID(userID uint) (models.Wallet, error) {
	wallet := models.Wallet{}
	result := w.db.Where("user_id = ?", userID).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return wallet, appErrors.ErrNotFound
	}
	return wallet, result.Error
}

func (w *WalletRepository) UpdateBalance(tx *gorm.DB, walletID uint, newBalance int) error {
	result := tx.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance)
	return result.Error
}

// row locking
func (w *WalletRepository) GetByUserIDForUpdate(tx *gorm.DB, userID uint) (models.Wallet, error) {
	wallet := models.Wallet{}
	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return wallet, appErrors.ErrNotFound
	}
	return wallet, result.Error
}
