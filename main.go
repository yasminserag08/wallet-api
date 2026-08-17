package main

import (
	"log"
	"wallet-api/handlers"
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
	authHandler := handlers.NewAuthHandler(userRepo)

	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.LogIn)

	router.Run()
}
