// Package response provides a standardised JSON envelope for all API responses.
//
// Every response follows the shape:
//
//	{
//	    "success": true|false,
//	    "message": "Human-readable string",
//	    "data":    <payload or null>,
//	    "errors":  <validation errors or null>,
//	    "meta":    <pagination/extra info or null>
//	}
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func Success(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, envelope{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, data any, message string) {
	c.JSON(http.StatusCreated, envelope{Success: true, Message: message, Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, statusCode int, message string, errors any) {
	c.JSON(statusCode, envelope{Success: false, Message: message, Errors: errors})
}

func ValidationError(c *gin.Context, errors any) {
	Error(c, http.StatusUnprocessableEntity, "Validation failed", errors)
}

func Unauthorized(c *gin.Context, message ...string) {
	msg := "Unauthenticated"
	if len(message) > 0 {
		msg = message[0]
	}
	Error(c, http.StatusUnauthorized, msg, nil)
}

func Forbidden(c *gin.Context) {
	Error(c, http.StatusForbidden, "You do not have permission to perform this action", nil)
}

func NotFound(c *gin.Context, resource string) {
	Error(c, http.StatusNotFound, resource+" not found", nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil)
}

func ServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "An unexpected error occurred", nil)
}

func TooManyRequests(c *gin.Context) {
	Error(c, http.StatusTooManyRequests, "Too many requests — please slow down", nil)
}

// Paginated wraps a list payload with pagination metadata.
func Paginated(c *gin.Context, data any, total, page, perPage int) {
	meta := gin.H{
		"total":    total,
		"page":     page,
		"per_page": perPage,
		"pages":    (total + perPage - 1) / perPage,
	}
	c.JSON(http.StatusOK, envelope{
		Success: true,
		Message: "OK",
		Data:    data,
		Meta:    meta,
	})
}
