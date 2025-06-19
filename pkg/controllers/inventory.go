package controllers

import (
	"errors"
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetInventories(c *gin.Context) {
	Inventorys, err := models.GetInventories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Inventorys)
}

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

func CreateInventory(c *gin.Context) {
	var inventory models.Inventory

	err := c.Bind(&inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
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

	if _, err := models.GetTenant(c, inventory.TenantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate tenant"})
		return
	}

	err = models.UpdateInventory(c, id, &inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetInventory(c, id)
	c.JSON(http.StatusOK, updated)
}