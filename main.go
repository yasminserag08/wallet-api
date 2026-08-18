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
	transactionHandler := handlers.NewTransactionHandler(transactionRepo, walletRepo)

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.LogIn)

	wallet := router.Group("/wallet")
	wallet.Use(middleware.RequireAuth())
	{
		wallet.GET("", walletHandler.GetWallet)
		wallet.POST("/deposit", walletHandler.Deposit)
		wallet.POST("/withdraw", walletHandler.Withdraw)
		wallet.POST("/transfer", walletHandler.Transfer)
	}

	transactions := router.Group("/transactions")
	transactions.Use(middleware.RequireAuth())
	{
		transactions.GET("", transactionHandler.ListTransactions)
	}

	router.Run()
}
