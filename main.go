package main

import (
	"log"
	"wallet-api/handlers"
	"wallet-api/middleware"
	"wallet-api/repositories"
	"wallet-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	router := gin.Default()

	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)
	walletService := services.NewWalletService(db, walletRepo, transactionRepo, userRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.LogIn)

	protected := router.Group("/wallet")
	protected.Use(middleware.RequireAuth())
	{
		protected.GET("", walletHandler.GetWallet)
		protected.POST("/deposit", walletHandler.Deposit)
		protected.POST("/withdraw", walletHandler.Withdraw)
		protected.POST("/transfer", walletHandler.Transfer)
	}

	router.Run()
}
