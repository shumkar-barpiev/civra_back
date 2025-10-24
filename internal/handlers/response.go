package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`  // "success" or "error"
	Data    interface{} `json:"data"`    // payload
	Message string      `json:"message"` // human-readable text
}

// Success sends a standardized success response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	})
}

// Created sends a standardized response for resource creation
func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Data:    nil,
		Message: message,
	})
}
