package utils

import (
	"errors"
	"gin-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StatusOK(c *gin.Context, data any, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

func StatusConflict(c *gin.Context, message any) {
	c.JSON(http.StatusConflict, models.Response{
		Status:  http.StatusConflict,
		Message: message,
	})
}

func StatusUnprocessable(c *gin.Context, message any) {
	c.JSON(http.StatusUnprocessableEntity, models.Response{
		Status:  http.StatusUnprocessableEntity,
		Message: message,
	})
}

func StatusBadRequest(c *gin.Context, message ...string) {
	msg := "bad request"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusBadRequest, models.Response{
		Status:  http.StatusBadRequest,
		Message: msg,
	})
}

func StatusServerError(c *gin.Context, message ...string) {
	msg := "internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(http.StatusInternalServerError, models.Response{
		Status:  http.StatusInternalServerError,
		Message: msg,
	})
}

func StatusNotFound(c *gin.Context, err error, message ...string) {
	msg := "the record not found"
	if len(message) > 0 {
		msg = message[0]
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  http.StatusNotFound,
			Message: msg,
		})
		return
	}

	StatusServerError(c)
}
