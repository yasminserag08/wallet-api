package handlers

import (
	"errors"
	"net/http"
	"wallet-api/services"

	"github.com/gin-gonic/gin"

	appErrors "wallet-api/errors"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// same structure used for both deposit and withdraw requests
type depositWithdrawRequest struct {
	Amount   int    `json:"amount" binding:"required,gt=0"` // greater than 0
	Category string `json:"category" binding:"required"`
	Note     string `json:"note"`
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	wallet, err := h.walletService.GetWallet(userID)

	if errors.Is(err, appErrors.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req depositWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.Deposit(userID, req.Amount, req.Category, req.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	// TODO
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	// TODO
}
