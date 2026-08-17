package services

import (
	"wallet-api/models"
	"wallet-api/repositories"

	appErrors "wallet-api/errors"

	"gorm.io/gorm"
)

type WalletService struct {
	db              *gorm.DB
	walletRepo      repositories.WalletRepositoryInterface
	transactionRepo repositories.TransactionRepositoryInterface
	userRepo        repositories.UserRepositoryInterface
}

func NewWalletService(db *gorm.DB, walletRepo repositories.WalletRepositoryInterface, transactionRepo repositories.TransactionRepositoryInterface, userRepo repositories.UserRepositoryInterface) *WalletService {
	return &WalletService{db: db, walletRepo: walletRepo, transactionRepo: transactionRepo, userRepo: userRepo}
}

func (s *WalletService) GetWallet(userID uint) (models.Wallet, error) {
	return s.walletRepo.GetByUserID(userID)
}

func (s *WalletService) Deposit(userID uint, amount int, category, note string) (models.Wallet, error) {
	var wallet models.Wallet

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		wallet, err = s.walletRepo.GetByUserIDForUpdate(tx, userID)
		if err != nil {
			return err
		}
		wallet.Balance += amount
		if err := s.walletRepo.UpdateBalance(tx, wallet.ID, wallet.Balance); err != nil {
			return err
		}

		return s.transactionRepo.Create(tx, models.Transaction{
			WalletID: wallet.ID,
			Type:     "deposit",
			Amount:   amount,
			Category: category,
			Note:     note,
		})
	})
	return wallet, err
}

func (s *WalletService) Withdraw(userID uint, amount int, category, note string) (models.Wallet, error) {
	var wallet models.Wallet

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		wallet, err = s.walletRepo.GetByUserIDForUpdate(tx, userID)
		if err != nil {
			return err
		}

		if wallet.Balance < amount {
			return appErrors.ErrInsufficientFunds
		}

		wallet.Balance -= amount
		if err := s.walletRepo.UpdateBalance(tx, wallet.ID, wallet.Balance); err != nil {
			return err
		}

		return s.transactionRepo.Create(tx, models.Transaction{
			WalletID: wallet.ID,
			Type:     "withdraw",
			Amount:   amount,
			Category: category,
			Note:     note,
		})
	})

	return wallet, err
}
