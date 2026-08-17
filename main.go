package main

import (
	"log"
	"wallet-api/handlers"
	"wallet-api/middleware"
	"wallet-api/repositories"

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
	authHandler := handlers.NewAuthHandler(userRepo)
	walletHandler := handlers.NewWalletHandler(walletRepo)

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.LogIn)

	protected := router.Group("/wallet")
	protected.Use(middleware.RequireAuth())
	{
		protected.GET("", walletHandler.GetWallet)
	}

	router.Run()
}
