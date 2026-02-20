package controllers

import (
	"net/http"
	"go-api/backend-api/database"
	"go-api/backend-api/helpers"
	"go-api/backend-api/models"
	"go-api/backend-api/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {

	// Inisialisasi struct for user login request and user model
	var req = structs.UserLoginRequest{}
	var user = models.User{}

	// Validasi input it's ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// find user by username, if not found return error Unauthorized
	// send respons error Unauthorized
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// compare password from request with hashed password in database, if not match return error Unauthorized
	// send respons error Unauthorized
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Invalid Credentials",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// generate token with username as payload
	token := helpers.GenerateToken(user.Username)

	// send response with user data and token
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
			Token:     &token,
		},
	})
}
