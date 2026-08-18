package services

import (
	"errors"
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

func (s *WalletService) Transfer(senderUserID uint, toUsername string, amount int, category, note string) error {
	// check if receiver and sender are the same
	receiver, err := s.userRepo.GetUserByUsername(toUsername)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return appErrors.ErrUserNotFound
		}
		return err
	}
	if receiver.ID == senderUserID {
		return appErrors.ErrSelfTransfer
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		sender, err := s.walletRepo.GetByUserIDForUpdate(tx, senderUserID)
		if err != nil {
			return err
		}

		if sender.Balance < amount {
			return appErrors.ErrInsufficientFunds
		}

		receiver, err := s.userRepo.GetUserByUsername(toUsername)

		receiverWallet, err := s.walletRepo.GetByUserIDForUpdate(tx, receiver.ID)
		if err != nil {
			return err
		}

		sender.Balance -= amount
		if err := s.walletRepo.UpdateBalance(tx, sender.ID, sender.Balance); err != nil {
			return err
		}

		receiverWallet.Balance += amount
		if err := s.walletRepo.UpdateBalance(tx, receiverWallet.ID, receiverWallet.Balance); err != nil {
			return err
		}

		if err := s.transactionRepo.Create(tx, models.Transaction{
			WalletID:        sender.ID,
			Type:            "transfer_out",
			Amount:          amount,
			Category:        category,
			Note:            note,
			RelatedWalletID: &receiverWallet.ID,
		}); err != nil {
			return err
		}

		return s.transactionRepo.Create(tx, models.Transaction{
			WalletID:        receiverWallet.ID,
			Type:            "transfer_in",
			Amount:          amount,
			Category:        category,
			Note:            note,
			RelatedWalletID: &sender.ID,
		})
	})
}
