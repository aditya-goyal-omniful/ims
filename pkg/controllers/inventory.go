package controllers

import (
	"errors"
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckInventoryRequest struct {
	SKUID    uuid.UUID `json:"sku_id" binding:"required"`
	HubID    uuid.UUID `json:"hub_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required"`
}


// GetInventories godoc
// @Summary Get all inventories
// @Tags Inventories
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {array} models.Inventory
// @Router /inventories [get]
func GetInventories(c *gin.Context) {
	Inventorys, err := models.GetInventories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Inventorys)
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

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inventory ID"})
		return
	}

	inventory, err := models.GetInventory(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
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

	err := c.Bind(&inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract tenant_id from header and assign to hub
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id in header"})
			return
		}
		inventory.TenantID = tenantID
	}

	if err := models.CreateInventory(c, &inventory); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory"})
		return
	}

	c.JSON(http.StatusCreated, inventory)
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
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inventory ID"})
		return
	}

	inventory, err := models.DeleteInventory(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inventory)
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

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inventory ID"})
		return
	}

	var inventory models.Inventory
	err = c.Bind(&inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if inventory.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, inventory.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate tenant"})
			return
		}
	}

	err = models.UpdateInventory(c, id, &inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetInventory(c, id)
	c.JSON(http.StatusOK, updated)
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

	err := c.Bind(&inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract tenant ID from header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id in header"})
		return
	}
	inventory.TenantID = tenantID

	// Call the upsert logic
	if err := models.UpsertInventory(c, &inventory); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory upserted", "inventory": inventory})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing tenant_id header or hub_id query param"})
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant_id"})
		return
	}

	hubID, err := uuid.Parse(hubIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hub_id"})
		return
	}

	view, err := models.GetInventoryWithDefaults(c, tenantID, hubID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, view)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	inv, err := models.GetInventoryBySkuHub(c, req.SKUID, req.HubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"available": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory"})
		return
	}

	if inv.Quantity < req.Quantity {
		c.JSON(http.StatusOK, gin.H{"available": false})
		return
	}

	// Deduct inventory
	newQty := inv.Quantity - req.Quantity
	if err := models.UpdateInventoryQuantity(c, inv.ID, newQty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": true})
}