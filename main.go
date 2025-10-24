package main

import (
	"log"
	controllers "zarrinpal/controller"
	"zarrinpal/db"
	"zarrinpal/middleware"

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
	router.POST("/login", controllers.LoginHandler())

	userAuthgroup := router.Group("/")
	userAuthgroup.Use(middleware.AuthMiddleware())
	{
		{
			userAuthgroup.POST("payment/request", controllers.RequestPaymentHandler())
			userAuthgroup.GET("payment/callback", controllers.CallbackHandler())
		}
	}

	log.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}
