package controllers

import (
	"go-api/backend-api/database"
	"go-api/backend-api/helpers"
	"go-api/backend-api/models"
	"go-api/backend-api/structs"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {
	var users []models.User
	var totalData int64
	
	// get pagination parameters from query string, with default values if not provided
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name, _ := c.GetQuery("name")

	db := database.DB.Model(&models.User{})

	// count total data without pagination
	db.Count(&totalData)

	// apply pagination scope and find users, if there's an error return response with status 500 Internal Server Error and error messages
	if err := db.Scopes(
		database.WithName(name),
		database.WithSort(c),
		database.Paginate(c)).Find(&users).Error; err != nil {
    c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
        Success: false,
        Message: "Failed to retrieve users",
        Errors: map[string]string{
            "error": err.Error(),
        },
    })
    return
}

	// compute total pages
	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	// send response with users data and pagination meta data
	c.JSON(http.StatusOK, structs.PaginatedResponse{
		Success: true,
		Message: "Successfully retrieved users",
		Data:    users,
		Pagination: structs.PaginationMeta{
			TotalData:   totalData,
			CurrentPage: page,
			TotalPages:  totalPages,
			Limit:       limit,
		},
	})
}

// func to FIndUserById
func FindUserById(c *gin.Context) {

	// get id from path parameter
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// find user by id, if there's an error return response with status 404 Not Found and error messages
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// send success response with user data
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved user",
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

// func to CreateUser
func CreateUser(c *gin.Context) {

	//struct user request
	var req = structs.UserCreateRequest{}

	// Bind JSON request ke struct UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Inisialisation struct User for create new user
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Save into database, if there's an error return response with status 500 Internal Server Error and error messages
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// success response with status 201 Created and user data
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

// func to UpdateUser patch user by id
func UpdateUser(c *gin.Context) {

	// get id from path parameter
	id := c.Param("id")

	// Inisialisation user
	var user models.User

	// check if user exists, if not return response with status 404 Not Found and error messages
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	//struct user request
	var req structs.UserUpdateRequest

	// Bind JSON request to struct UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Update user fields if they are provided in the request
	if req.Password != "" {
    hashedPassword := helpers.HashPassword(req.Password)
    
    req.Password = hashedPassword
	}
	// save updated user into database, if there's an error return response with status 500 Internal Server Error and error messages
	if err := database.DB.Model(&user).Updates(req).Error; err != nil {
      c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
      	Success: false,
      	Message: "Failed to update user",
				Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// success response with status 200 OK and updated user data
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User updated successfully",
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

func DeleteUser(c *gin.Context) {

	// get id from path parameter
	id := c.Param("id")

	// Inisialisasi user
	var user models.User

	// check if user exists, if not return response with status 404 Not Found and error messages
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Delete user from database, if there's an error return response with status 500 Internal Server Error and error messages
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// success response with status 200 OK and message user deleted successfully
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

