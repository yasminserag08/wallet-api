package handlers

import (
	"net/http"
	"strconv"
	"time"
	"wallet-api/repositories"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionRepo repositories.TransactionRepositoryInterface
	walletRepo      repositories.WalletRepositoryInterface
}

func NewTransactionHandler(transactionRepo repositories.TransactionRepositoryInterface, walletRepo repositories.WalletRepositoryInterface) *TransactionHandler {
	return &TransactionHandler{transactionRepo: transactionRepo, walletRepo: walletRepo}
}

func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	wallet, err := h.walletRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	// default is page 1 and limit 10 if not specified
	filter := repositories.TransactionFilter{
		Category: c.Query("category"),
		Page:     1,
		Limit:    10,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.From = t
		}
	}

	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			filter.To = t
		}
	}

	transactions, err := h.transactionRepo.GetByWalletID(wallet.ID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
