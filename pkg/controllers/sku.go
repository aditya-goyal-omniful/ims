package controllers

import (
	"context"
	"errors"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"gorm.io/gorm"
)

// GetSkus

type SkuFetcher interface {
	GetFilteredSkus(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error)
}

func getSkusLogic(service SkuFetcher, tenantIDStr, sellerIDStr string, skuCodes []string) ([]models.Sku, int, error) {
	if tenantIDStr == "" {
		return nil, int(http.StatusBadRequest), errors.New("missing X-Tenant-ID header")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid tenant_id in header")
	}

	var sellerID uuid.UUID
	if sellerIDStr != "" {
		sellerID, err = uuid.Parse(sellerIDStr)
		if err != nil {
			return nil, int(http.StatusBadRequest), errors.New("invalid seller_id")
		}
	}

	skus, err := service.GetFilteredSkus(context.Background(), tenantID, sellerID, skuCodes)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return skus, int(http.StatusOK), nil
}

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
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	sellerIDStr := c.Query("seller_id")
	skuCodes := c.QueryArray("sku_codes")

	skus, status, err := getSkusLogic(models.SKUModel{}, tenantIDStr, sellerIDStr, skuCodes)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, skus)
}

// GetSkuByID

type SkuGetter interface {
	GetSku(ctx context.Context, id uuid.UUID) (*models.Sku, error)
}

func getSkuByIDLogic(service SkuGetter, idStr string) (*models.Sku, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid SKU ID")
	}

	sku, err := service.GetSku(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return sku, int(http.StatusOK), nil
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

	sku, status, err := getSkuByIDLogic(models.SKUModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, sku)
}

// CreateSku

type SkuCreator interface {
	CreateSku(ctx context.Context, sku *models.Sku) error
}

func createSkuLogic(service SkuCreator, tenantIDStr string, sku *models.Sku) (int, error) {
	if tenantIDStr == "" {
		return int(http.StatusBadRequest), errors.New("missing X-Tenant-ID header")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return int(http.StatusBadRequest), errors.New("invalid tenant_id in header")
	}
	sku.TenantID = tenantID

	if err := service.CreateSku(context.Background(), sku); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return int(http.StatusBadRequest), errors.New("tenant or seller not found")
		}
		return int(http.StatusInternalServerError), errors.New("failed to create sku")
	}

	return int(http.StatusCreated), nil
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

	if err := c.Bind(&sku); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	status, err := createSkuLogic(models.SKUModel{}, c.GetHeader("X-Tenant-ID"), &sku)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, sku)
}

// DeleteSku

type SkuDeleter interface {
	DeleteSku(ctx context.Context, id uuid.UUID) (*models.Sku, error)
}

func deleteSkuLogic(service SkuDeleter, idStr string) (*models.Sku, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid sku id")
	}

	sku, err := service.DeleteSku(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusNotFound), errors.New("sku not found")
	}

	return sku, int(http.StatusOK), nil
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

	sku, status, err := deleteSkuLogic(models.SKUModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, sku)
}

// UpdateSku

type SkuUpdater interface {
	UpdateSku(ctx context.Context, id uuid.UUID, updated *models.Sku) error
	GetSku(ctx context.Context, id uuid.UUID) (*models.Sku, error)
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func updateSkuLogic(service SkuUpdater, idStr string, sku *models.Sku) (*models.Sku, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid sku id")
	}

	if sku.TenantID != uuid.Nil {
		if _, err := service.GetTenant(context.Background(), sku.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, int(http.StatusBadRequest), errors.New("tenant not found")
			}
			return nil, int(http.StatusInternalServerError), errors.New("failed to validate tenant")
		}
	}

	if err := service.UpdateSku(context.Background(), id, sku); err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	updated, err := service.GetSku(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return updated, int(http.StatusOK), nil
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

	var sku models.Sku
	if err := c.Bind(&sku); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	updated, status, err := updateSkuLogic(models.SKUModel{}, idStr, &sku)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, updated)
}
