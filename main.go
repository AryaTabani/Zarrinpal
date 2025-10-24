package main

import (
	"log"
	controllers "zarrinpal/controller"
	"zarrinpal/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables.")
	}

	db.InitDB()

	router := gin.Default()
	router.POST("/register", controllers.UserRegisterHandler())
	paymentRoutes := router.Group("/payment")
	{
		paymentRoutes.POST("/request", controllers.RequestPaymentHandler())
		paymentRoutes.GET("/callback", controllers.CallbackHandler())
	}

	log.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
