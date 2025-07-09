package controllers

import (
	"context"
	"errors"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

type Validator interface {
	ValidateHubAndSku(ctx context.Context, hubID, skuID uuid.UUID) (bool, error)
}

func validateOrderLogic(service Validator, hubIDStr, skuIDStr string) (bool, int, error) {
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		return false, int(http.StatusBadRequest), errors.New("invalid hub_id")
	}

	skuID, err := uuid.Parse(skuIDStr)
	if err != nil {
		return false, int(http.StatusBadRequest), errors.New("invalid sku_id")
	}

	isValid, err := service.ValidateHubAndSku(context.Background(), hubID, skuID)
	if err != nil {
		return false, int(http.StatusInternalServerError), errors.New("validation failed")
	}

	return isValid, int(http.StatusOK), nil
}

// ValidateOrder godoc
// @Summary Validate hub and SKU IDs
// @Tags Validators
// @Produce json
// @Param hub_id path string true "Hub ID"
// @Param sku_id path string true "SKU ID"
// @Success 200 {object} map[string]bool
// @Router /validators/validate_order/{hub_id}/{sku_id} [get]
func ValidateOrder(c *gin.Context) {
	hubIDStr := c.Param("hub_id")
	skuIDStr := c.Param("sku_id")

	log.Infof(i18n.Translate(c, "Received validate request for hubID=%s skuID=%s"), hubIDStr, skuIDStr)

	isValid, status, err := validateOrderLogic(models.InventoryModel{}, hubIDStr, skuIDStr)
	if err != nil {
		log.Errorf("Validation failed: %v", err)
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, gin.H{i18n.Translate(c, "is_valid"): isValid})
}
