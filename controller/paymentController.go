package controllers

import (
	"net/http"
	"zarrinpal/models"
	services "zarrinpal/service"

	"github.com/gin-gonic/gin"
)

func RequestPaymentHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.PaymentRequestPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse[any]{
				Success: false,
				Error:   "Invalid request body: " + err.Error(),
			})
			return
		}

		paymentURL, err := services.RequestPayment(c.Request.Context(), &payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse[any]{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse[models.PaymentURLResponse]{
			Success: true,
			Data:    models.PaymentURLResponse{PaymentURL: paymentURL},
		})
	}
}

func CallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authority := c.Query("Authority")
		status := c.Query("Status")

		if authority == "" || status == "" {
			c.String(http.StatusBadRequest, "Invalid callback parameters from Zarinpal.")
			return
		}

		refID, err := services.VerifyPayment(c.Request.Context(), authority, status)
		if err != nil {
			c.String(http.StatusBadRequest, "Payment verification failed: "+err.Error())
			return
		}

		c.String(http.StatusOK, "Payment successful! Reference ID: "+refID)
	}
}
