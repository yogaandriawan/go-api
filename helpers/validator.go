package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10" 
	"gorm.io/gorm"
)

// TranslateErrorMessage for translating validation and database error messages into user-friendly format
func TranslateErrorMessage(err error) map[string]string {
	// create a map to hold error messages, key is the field name and value is the error message
	errorsMap := make(map[string]string)

	// Handle validasi from validator.v10
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field() // field name errors
			switch fieldError.Tag() {   // tag validasi error
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field) // return error message if field is required but not provided
			case "email":
				errorsMap[field] = "Invalid email format" // return error message if email format is invalid
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field) // return error message if data already exists
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param()) // return error message if value is too short
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param()) // return error message if value is too long
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field) // return error message if value is not a number
			default:
				errorsMap[field] = "Invalid value" // return error message default for other validation errors
			}
		}
	}

	// Handle error form GORM untuk duplicate entry
	if err != nil {
		// check if error message contains "Duplicate entry" and field name, then return appropriate error message for each field
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				errorsMap["Username"] = "Username already exists" // return error message if username already exists
			}
			if strings.Contains(err.Error(), "email") {
				errorsMap["Email"] = "Email already exists" // return error message if email already exists
			}
		} else if err == gorm.ErrRecordNotFound {
			// handle error for record not found, you can customize the message based on your needs, here I just return "Record not found"
			errorsMap["Error"] = "Record not found"
		}
	}

	// return the map of error messages
	return errorsMap
}

// IsDuplicateEntryError to check if the error is a duplicate entry error from GORM
func IsDuplicateEntryError(err error) bool {
	// check if error message contains "Duplicate entry"
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
