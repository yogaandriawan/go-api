package database

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//get pagination parameters from query string, with default values if not provided
		page, _ := strconv.Atoi(c.Query("page"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.Query("limit"))
		switch {
		case limit > 100:
			limit = 100 // Maksimal limit is 100 to prevent abuse
		case limit <= 0:
			limit = 10
		}

		// calculate offset for pagination
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func WithSort(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sort := c.DefaultQuery("sort", "created_at") // Default sort by created_at
		order := c.DefaultQuery("order", "desc")    // Default order desc

		// it's should be 'asc' or 'desc'
		if order != "asc" && order != "desc" {
			order = "desc"
		}

		query := fmt.Sprintf("%s %s", sort, order)
		return db.Order(query)
	}
}

// Filter a names (Partial Match)
func WithName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if name == "" {
			return db
		}
		return db.Where("name ILIKE ?", "%"+name+"%")
	}
}

//Filter for status (Exact Match)
func WithStatus(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == "" {
			return db
		}
		return db.Where("status = ?", status)
	}
}

// Filter for created_at between start and end date (Range Filter)
func CreatedBetween(start, end string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start == "" || end == "" {
			return db
		}
		return db.Where("created_at BETWEEN ? AND ?", start, end)
	}
}