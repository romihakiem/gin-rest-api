package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func GetUserID(c *gin.Context) uint {
	userId, exists := c.Get("userAuth")
	if !exists {
		return 0
	}
	return userId.(uint)
}

func IsExist(db *gorm.DB, table, field string, value any) bool {
	var count int64

	result := db.Table(table).Where(field+" = ?", value).Count(&count)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	}

	return count > 0
}

func FormatErrors(errs validator.ValidationErrors) map[string]string {
	message := make(map[string]string)

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			message[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			message[err.Field()] = fmt.Sprintf("%s must be a valid email address", err.Field())
		case "min":
			message[err.Field()] = fmt.Sprintf("%s must have at least %s characters", err.Field(), err.Param())
		case "max":
			message[err.Field()] = fmt.Sprintf("%s must have at most %s characters", err.Field(), err.Param())
		case "gt":
			message[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "gte":
			message[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
		default:
			message[err.Field()] = fmt.Sprintf("Validation validations on field %s", err.Field())
		}
	}

	return message
}
