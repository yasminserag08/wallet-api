package repositories

import (
	"time"
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

type TransactionFilter struct {
	Category string
	From     time.Time
	To       time.Time
	Page     int
	Limit    int
}

func (r *TransactionRepository) GetByWalletID(walletID uint, filter TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.db.Where("wallet_id = ?", walletID)

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if !filter.From.IsZero() {
		query = query.Where("created_at >= ?", filter.From)
	}

	if !filter.To.IsZero() {
		query = query.Where("created_at <= ?", filter.To)
	}

	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit).Order("created_at desc")

	result := query.Find(&transactions)
	return transactions, result.Error
}

type CategorySummary struct {
	Category string `json:"category"`
	Total    int    `json:"total"`
}

func (r *TransactionRepository) GetSummary(walletID uint) ([]CategorySummary, error) {
	var summary []CategorySummary

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result := r.db.Model(&models.Transaction{}).
		Select("category, SUM(amount) as total").
		Where("wallet_id = ? AND created_at >= ? AND type IN ?", walletID, startOfMonth, []string{"withdraw", "transfer_out"}).Group("category").
		Scan(&summary)
		// only show expenses (withdrawal/transfer out) in the summary

	return summary, result.Error
}
