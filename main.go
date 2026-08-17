package main

import (
	"log"

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

	router.GET("/wallet/hello", func(c *gin.Context) { c.JSON(200, gin.H{"message": "hello"}) })

	router.Run()
}
