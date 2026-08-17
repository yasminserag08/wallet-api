package repositories

import (
	"wallet-api/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) Create(tx *gorm.DB, transaction models.Transaction) error {
	return tx.Create(&transaction).Error
}
