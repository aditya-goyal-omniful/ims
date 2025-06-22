package controllers

import (
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
)

func ValidateOrder(c *gin.Context) {
	log.Infof("Received validate request for hubID=%s skuID=%s", c.Param("hub_id"), c.Param("sku_id"))

	hubIDStr := c.Param("hub_id")
	skuIDStr := c.Param("sku_id")

	hubID, err1 := uuid.Parse(hubIDStr)
	skuID, err2 := uuid.Parse(skuIDStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hub_id or sku_id"})
		return
	}

	isValid, err := models.ValidateHubAndSku(c, hubID, skuID)
	if err != nil {
		log.Errorf("Validation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_valid": isValid})
}