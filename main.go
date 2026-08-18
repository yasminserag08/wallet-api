package main

import (
	"log"
	"wallet-api/db"
	"wallet-api/handlers"
	"wallet-api/middleware"
	"wallet-api/repositories"
	"wallet-api/services"

	_ "wallet-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Wallet API
// @version         1.0
// @description     A mini wallet and expense tracker API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and JWT token.
func main() {
	db, err := db.Connect()
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
		transactions.GET("/summary", transactionHandler.GetSummary)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
