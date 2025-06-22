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

	// Parse UUIDs
	hubID, err1 := uuid.Parse(hubIDStr)
	skuID, err2 := uuid.Parse(skuIDStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hub_id, sku_id"})
		return
	}

	// Validate each individually using model layer
	hub, err := models.GetHub(c, hubID)
	if err != nil || hub == nil {
		c.JSON(http.StatusOK, gin.H{"is_valid": false})
		return
	}

	sku, err := models.GetSku(c, skuID)
	if err != nil || sku == nil {
		c.JSON(http.StatusOK, gin.H{"is_valid": false})
		return
	}

	// All validations passed
	c.JSON(http.StatusOK, gin.H{"is_valid": true})
}