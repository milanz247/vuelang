package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-cloud-erp/app/models"
)

// GreetingController handles requests for /api/v1/greeting.
type GreetingController struct{}

// Index  GET /api/v1/greeting
func (ctrl *GreetingController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewGreeting())
}
