package controllers

import (
	"fmt"
	"net/http"
	"go-api/backend-api/database"
	"go-api/backend-api/helpers"
	"go-api/backend-api/models"
	"go-api/backend-api/structs"

	"github.com/gin-gonic/gin"
)

// Register for handling user registration, it will receive JSON request with user data, validate it, create a new user in the database, and return the created user data as response
func Register(c *gin.Context) {
	// Inisialisasi struct for user creation request
	var req = structs.UserCreateRequest{}
	fmt.Println("Received registration request: ", req)


	// Validasi request should be a JSON 
	if err := c.ShouldBindJSON(&req); err != nil {
		// if validation error, return response with status 422 Unprocessable Entity and error messages
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validasi Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Buat instance user baru dengan data dari request, pastikan untuk hash password sebelum menyimpan ke database
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// save in database, if there's an error, check if it's a duplicate entry error and return appropriate response
	if err := database.DB.Create(&user).Error; err != nil {
		// Check if the error is a duplicate entry error, if so return 409 Conflict, otherwise return 500 Internal Server Error
		if helpers.IsDuplicateEntryError(err) {
			// if it's a duplicate entry error, return response with status 409 Conflict and error messages
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Success: false,
				Message: "Duplicate entry error",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		} else {
			// if it's other error, return response with status 500 Internal Server Error and error messages
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Failed to create user",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		}
		return
	}

	// return response with status 201 Created and the created user data
	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
