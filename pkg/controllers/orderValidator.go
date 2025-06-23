package controllers

import (
	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

// ValidateOrder godoc
// @Summary Validate hub and SKU IDs
// @Tags Validators
// @Produce json
// @Param hub_id path string true "Hub ID"
// @Param sku_id path string true "SKU ID"
// @Success 200 {object} map[string]bool
// @Router /validators/validate_order/{hub_id}/{sku_id} [get]
func ValidateOrder(c *gin.Context) {
	log.Infof(i18n.Translate(c, "Received validate request for hubID=%s skuID=%s"), c.Param("hub_id"), c.Param("sku_id"))

	hubIDStr := c.Param("hub_id")
	skuIDStr := c.Param("sku_id")

	hubID, err1 := uuid.Parse(hubIDStr)
	skuID, err2 := uuid.Parse(skuIDStr)

	if err1 != nil || err2 != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid hub_id or sku_id")})
		return
	}

	isValid, err := models.ValidateHubAndSku(c, hubID, skuID)
	if err != nil {
		log.Errorf("Validation failed: %v", err)
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Internal server error")})
		return
	}

	c.JSON(int(http.StatusOK), gin.H{i18n.Translate(c, "is_valid"): isValid})
}