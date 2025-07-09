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

type CheckInventoryRequest struct {
	SKUID    uuid.UUID `json:"sku_id" binding:"required"`
	HubID    uuid.UUID `json:"hub_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required"`
}

// GetInventories

type InventoryFetcher interface {
	GetInventories(ctx context.Context) ([]models.Inventory, error)
}

func getInventoriesLogic(service InventoryFetcher) ([]models.Inventory, int, error) {
	inventories, err := service.GetInventories(context.Background())
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}
	return inventories, int(http.StatusOK), nil
}

// GetInventories godoc
// @Summary Get all inventories
// @Tags Inventories
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {array} models.Inventory
// @Router /inventories [get]
func GetInventories(c *gin.Context) {
	inventories, status, err := getInventoriesLogic(models.InventoryModel{})
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}
	c.JSON(status, inventories)
}

// GetInventoryByID

type InventoryByIDFetcher interface {
	GetInventory(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
}

func getInventoryByIDLogic(service InventoryByIDFetcher, idStr string) (*models.Inventory, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid inventory id")
	}

	inventory, err := service.GetInventory(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return inventory, int(http.StatusOK), nil
}

// GetInventoryByID godoc
// @Summary Get inventory by ID
// @Tags Inventories
// @Produce json
// @Param id path string true "Inventory ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Inventory
// @Router /inventories/{id} [get]
func GetInventoryByID(c *gin.Context) {
	idStr := c.Param("id")

	inventory, status, err := getInventoryByIDLogic(models.InventoryModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, inventory)
}

// CreateInventory

type InventoryCreator interface {
	CreateInventory(ctx context.Context, inv *models.Inventory) error
}

func createInventoryLogic(service InventoryCreator, tenantIDStr string, inventory *models.Inventory) (int, error) {
	// Validate tenant ID
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return int(http.StatusBadRequest), errors.New("invalid tenant_id in header")
		}
		inventory.TenantID = tenantID
	}

	// Save to DB
	if err := service.CreateInventory(context.Background(), inventory); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return int(http.StatusBadRequest), errors.New("tenant not found")
		}
		return int(http.StatusInternalServerError), errors.New("failed to create inventory")
	}

	return int(http.StatusCreated), nil
}

// CreateInventory godoc
// @Summary Create new inventory
// @Tags Inventories
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param inventory body models.Inventory true "Inventory to create"
// @Success 201 {object} models.Inventory
// @Router /inventories [post]
func CreateInventory(c *gin.Context) {
	var inventory models.Inventory

	if err := c.Bind(&inventory); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	status, err := createInventoryLogic(models.InventoryModel{}, c.GetHeader("X-Tenant-ID"), &inventory)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, inventory)
}

// DeleteInventory

type InventoryDeleter interface {
	DeleteInventory(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
}

func deleteInventoryLogic(service InventoryDeleter, idStr string) (*models.Inventory, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid inventory id")
	}

	inv, err := service.DeleteInventory(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, int(http.StatusNotFound), errors.New("inventory not found")
		}
		return nil, int(http.StatusInternalServerError), err
	}

	return inv, int(http.StatusOK), nil
}

// DeleteInventory godoc
// @Summary Delete inventory by ID
// @Tags Inventories
// @Produce json
// @Param id path string true "Inventory ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Inventory
// @Router /inventories/{id} [delete]
func DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")

	inv, status, err := deleteInventoryLogic(models.InventoryModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, inv)
}

// UpdateInventory

type InventoryUpdater interface {
	UpdateInventory(ctx context.Context, id uuid.UUID, updated *models.Inventory) error
	GetInventory(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
}

type TenantValidator interface {
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func updateInventoryLogic(
	service InventoryUpdater,
	tenantService TenantValidator,
	idStr string,
	inventory *models.Inventory,
) (*models.Inventory, int, error) {
	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid inventory id")
	}

	// Tenant validation (optional)
	if inventory.TenantID != uuid.Nil {
		_, err := tenantService.GetTenant(context.Background(), inventory.TenantID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, int(http.StatusBadRequest), errors.New("tenant not found")
			}
			return nil, int(http.StatusInternalServerError), errors.New("failed to validate tenant")
		}
	}

	// Update inventory
	if err := service.UpdateInventory(context.Background(), id, inventory); err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	// Fetch updated
	updated, err := service.GetInventory(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return updated, int(http.StatusOK), nil
}

// UpdateInventory godoc
// @Summary Update inventory by ID
// @Tags Inventories
// @Accept json
// @Produce json
// @Param id path string true "Inventory ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param inventory body models.Inventory true "Updated inventory"
// @Success 200 {object} models.Inventory
// @Router /inventories/{id} [put]
func UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")

	var inventory models.Inventory
	if err := c.Bind(&inventory); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	updated, status, err := updateInventoryLogic(
		models.InventoryModel{},
		models.TenantModel{},
		idStr,
		&inventory,
	)

	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, updated)
}

// UpsertInventory

type InventoryUpserter interface {
	UpsertInventory(ctx context.Context, inv *models.Inventory) error
}

func upsertInventoryLogic(service InventoryUpserter, tenantIDStr string, inv *models.Inventory) (int, error) {
	// Parse tenant ID
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return int(http.StatusBadRequest), errors.New("invalid tenant_id in header")
	}
	inv.TenantID = tenantID

	// Call DB upsert
	if err := service.UpsertInventory(context.Background(), inv); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return int(http.StatusBadRequest), errors.New("tenant not found")
		}
		return int(http.StatusInternalServerError), errors.New("failed to upsert inventory")
	}

	return int(http.StatusOK), nil
}

// UpsertInventory godoc
// @Summary Upsert (create or update) inventory
// @Tags Inventories
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param inventory body models.Inventory true "Inventory object"
// @Success 200 {object} map[string]string
// @Router /inventories/upsert [post]
func UpsertInventory(c *gin.Context) {
	var inventory models.Inventory

	if err := c.Bind(&inventory); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	status, err := upsertInventoryLogic(models.InventoryModel{}, c.GetHeader("X-Tenant-ID"), &inventory)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, gin.H{
		i18n.Translate(c, "message"):    i18n.Translate(c, "Inventory upserted"),
		i18n.Translate(c, "inventory"):  inventory,
	})
}

// ViewInventoryWithDefaults

type InventoryViewer interface {
	GetInventoryWithDefaults(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error)
}

func viewInventoryWithDefaultsLogic(service InventoryViewer, tenantIDStr, hubIDStr string) ([]models.InventoryView, int, error) {
	// Validate tenant UUID
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid tenant_id")
	}

	// Validate hub UUID
	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid hub_id")
	}

	// Fetch inventory view
	result, err := service.GetInventoryWithDefaults(context.Background(), tenantID, hubID)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return result, int(http.StatusOK), nil
}

// ViewInventoryWithDefaults godoc
// @Summary View inventory including SKUs with zero quantity
// @Tags Inventories
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param hub_id query string true "Hub ID"
// @Success 200 {array} models.InventoryView
// @Router /inventories/view [get]
func ViewInventoryWithDefaults(c *gin.Context) {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	hubIDStr := c.Query("hub_id")

	if tenantIDStr == "" || hubIDStr == "" {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Missing tenant_id header or hub_id query param")})
		return
	}

	view, status, err := viewInventoryWithDefaultsLogic(models.InventoryModel{}, tenantIDStr, hubIDStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, view)
}

// CheckAndUpdateInventory

type InventoryChecker interface {
	GetInventoryBySkuHub(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error)
	UpdateInventoryQuantity(ctx context.Context, invID uuid.UUID, newQty int) error
}

func checkAndUpdateInventoryLogic(service InventoryChecker, req CheckInventoryRequest) (bool, int, error) {
	inv, err := service.GetInventoryBySkuHub(context.Background(), req.SKUID, req.HubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, int(http.StatusOK), nil // not available, but valid
		}
		return false, int(http.StatusInternalServerError), errors.New("failed to fetch inventory")
	}

	if inv.Quantity < req.Quantity {
		return false, int(http.StatusOK), nil
	}

	newQty := inv.Quantity - req.Quantity
	if err := service.UpdateInventoryQuantity(context.Background(), inv.ID, newQty); err != nil {
		return false, int(http.StatusInternalServerError), errors.New("failed to update inventory")
	}

	return true, int(http.StatusOK), nil
}


// CheckAndUpdateInventory godoc
// @Summary Check and update inventory if sufficient
// @Tags Inventories
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param payload body CheckInventoryRequest true "Inventory check payload"
// @Success 200 {object} map[string]bool
// @Router /inventory/check-and-update [post]
func CheckAndUpdateInventory(c *gin.Context) {
	var req CheckInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request payload")})
		return
	}

	available, status, err := checkAndUpdateInventoryLogic(models.InventoryModel{}, req)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, gin.H{i18n.Translate(c, "available"): available})
}
