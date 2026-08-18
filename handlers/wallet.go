package handlers

import (
	"errors"
	"net/http"
	"strconv"
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
type DepositWithdrawRequest struct {
	Amount   int    `json:"amount" binding:"required,gt=0"` // greater than 0
	Category string `json:"category" binding:"required"`
	Note     string `json:"note"`
}

type TransferRequest struct {
	ToUsername string `json:"toUsername" binding:"required"`
	Amount     int    `json:"amount" binding:"required,gt=0"`
	Category   string `json:"category" binding:"required"`
	Note       string `json:"note"`
}

// @Summary      Get wallet
// @Description  Get current user's wallet and balance
// @Tags         wallet
// @Produce      json
// @Success      200  {object}  models.Wallet
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	role := c.MustGet("role").(string)

	if role == "admin" {
		if queryID := c.Query("userID"); queryID != "" {
			id, err := strconv.ParseUint(queryID, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
				return
			}
			userID = uint(id)
		}
	}

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

// @Summary      Deposit
// @Description  Deposit money into current user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        request body DepositWithdrawRequest true "Deposit request"
// @Success      200  {object}  models.Wallet
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /wallet/deposit [post]
func (h *WalletHandler) Deposit(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req DepositWithdrawRequest
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

// @Summary      Withdraw
// @Description  Withdraw money from current user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        request body DepositWithdrawRequest true "Withdraw request"
// @Success      200  {object}  models.Wallet
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /wallet/withdraw [post]
func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req DepositWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.Withdraw(userID, req.Amount, req.Category, req.Note)
	if errors.Is(err, appErrors.ErrInsufficientFunds) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// @Summary      Transfer
// @Description  Transfer money to another user
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        request body TransferRequest true "Transfer request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /wallet/transfer [post]
func (h *WalletHandler) Transfer(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.walletService.Transfer(userID, req.ToUsername, req.Amount, req.Category, req.Note); err != nil {
		if errors.Is(err, appErrors.ErrInsufficientFunds) || errors.Is(err, appErrors.ErrSelfTransfer) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, appErrors.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipient not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}
