package controllers

import (
	"go-api/backend-api/database"
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
