package main

import (
	"example/bookstore/database"
	"example/bookstore/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := gin.Default()
	routes.RegisterBookRoutes(router)
	router.Run("localhost:8080")
}
