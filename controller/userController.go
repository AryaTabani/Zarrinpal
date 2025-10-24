package controllers

import (
	"errors"
	"net/http"
	"zarrinpal/models"
	services "zarrinpal/service"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.RegisterPayload
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			response := models.APIResponse[any]{
				Success: false,
				Error:   "invalid request body: " + err.Error(),
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		_, err = services.RegisterUser(c.Request.Context(), &payload)
		if err != nil {
			if errors.Is(err, services.ErrUserExists) {
				response := models.APIResponse[any]{
					Success: false,
					Error:   err.Error(),
				}
				c.JSON(http.StatusConflict, response)
				return
			}
			response := models.APIResponse[any]{
				Success: false,
				Error:   "Failed to create user" + err.Error(),
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response := models.APIResponse[any]{
			Success: true,
			Message: "user created successfully",
		}
		c.JSON(http.StatusOK, response)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.LoginPayload
		err := c.ShouldBindBodyWithJSON(payload)
		if err != nil {
			response := models.APIResponse[any]{
				Success: false,
				Error:   "invalid request body",
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		token, err := services.LoginUser(c.Request.Context(), &payload)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				response := models.APIResponse[any]{
					Success: false,
					Error:   err.Error(),
				}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
			response := models.APIResponse[any]{
				Success: false,
				Error:   "Login Failed",
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		response := models.APIResponse[models.LoginResponse]{
			Success: true,
			Data:    models.LoginResponse{Token: token},
		}
		c.JSON(http.StatusOK, response)

	}
}
