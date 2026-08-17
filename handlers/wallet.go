package handlers

import (
	"errors"
	"net/http"
	"wallet-api/repositories"

	"github.com/gin-gonic/gin"

	appErrors "wallet-api/errors"
)

type WalletHandler struct {
	walletRepo repositories.WalletRepositoryInterface
}

func NewWalletHandler(walletRepo repositories.WalletRepositoryInterface) *WalletHandler {
	return &WalletHandler{walletRepo: walletRepo}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	wallet, err := h.walletRepo.GetByUserID(userID)

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
