package controllers

import (
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	Inventory, err := models.GetInventory(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Inventory)
}

func CreateInventory(c *gin.Context) {
	var Inventory models.Inventory

	err := c.Bind(&Inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.CreateInventory(c, &Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Inventory)
}

func DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inventory ID"})
		return
	}

	Inventory, err := models.DeleteInventory(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, Inventory)
}

func UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inventory ID"})
		return
	}

	var Inventory models.Inventory
	err = c.Bind(&Inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.UpdateInventory(c, id, &Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetInventory(c, id)
	c.JSON(http.StatusOK, updated)
}