package controllers

import (
	"errors"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"gorm.io/gorm"
)

// GetSkus godoc
// @Summary Get all SKUs (with optional filters)
// @Tags SKUs
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param seller_id query string false "Filter by Seller ID"
// @Param sku_codes query []string false "Filter by multiple SKU codes (repeat param)"
// @Success 200 {array} models.Sku
// @Router /skus [get]
func GetSkus(c *gin.Context) {
	// Extract tenant ID from auth headers
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Missing X-Tenant-ID header")})
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid tenant_id in header")})
		return
	}

	// Optional filters
	sellerIDStr := c.Query("seller_id")
	skuCodes := c.QueryArray("sku_codes") // Supports ?sku_codes=abc&sku_codes=def

	var sellerID uuid.UUID
	if sellerIDStr != "" {
		sellerID, err = uuid.Parse(sellerIDStr)
		if err != nil {
			c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid seller_id")})
			return
		}
	}

	skus, err := models.GetFilteredSkus(c, tenantID, sellerID, skuCodes)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), skus)
}

// GetSkuByID godoc
// @Summary Get SKU by ID
// @Tags SKUs
// @Produce json
// @Param id path string true "SKU ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Sku
// @Router /skus/{id} [get]
func GetSkuByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Sku ID")})
		return
	}

	Sku, err := models.GetSku(c, id)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), Sku)
}

// CreateSku godoc
// @Summary Create a new SKU
// @Tags SKUs
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sku body models.Sku true "SKU to create"
// @Success 201 {object} models.Sku
// @Router /skus [post]
func CreateSku(c *gin.Context) {
	var sku models.Sku

	err := c.Bind(&sku)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	// Extract tenant_id from header and assign to hub
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid tenant_id in header")})
			return
		}
		sku.TenantID = tenantID
	}

	if err := models.CreateSku(c, &sku); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant or Seller not found")})
			return
		}
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to create sku")})
		return
	}

	c.JSON(int(http.StatusCreated), sku)
}

// DeleteSku godoc
// @Summary Delete SKU by ID
// @Tags SKUs
// @Produce json
// @Param id path string true "SKU ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Sku
// @Router /skus/{id} [delete]
func DeleteSku(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Sku ID")})
		return
	}

	sku, err := models.DeleteSku(c, id)
	if err != nil {
		c.JSON(int(http.StatusNotFound), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Sku not found")})
		return
	}

	c.JSON(int(http.StatusOK), sku)
}

// UpdateSku godoc
// @Summary Update SKU by ID
// @Tags SKUs
// @Accept json
// @Produce json
// @Param id path string true "SKU ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sku body models.Sku true "Updated SKU"
// @Success 200 {object} models.Sku
// @Router /skus/{id} [put]
func UpdateSku(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Sku ID")})
		return
	}

	var sku models.Sku
	err = c.Bind(&sku)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	if sku.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, sku.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
				return
			}
			c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to validate tenant")})
			return
		}
	}

	err = models.UpdateSku(c, id, &sku)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	updated, _ := models.GetSku(c, id)
	c.JSON(int(http.StatusOK), updated)
}