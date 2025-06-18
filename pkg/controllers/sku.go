package controllers

import (
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSkus(c *gin.Context) {
	skus, err := models.GetInventories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skus)
}

func GetSkuByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sku ID"})
		return
	}

	Sku, err := models.GetSku(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Sku)
}

func CreateSku(c *gin.Context) {
	var sku models.Sku

	err := c.Bind(&sku)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.CreateSku(c, &sku)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sku)
}

func DeleteSku(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sku ID"})
		return
	}

	sku, err := models.DeleteSku(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sku not found"})
		return
	}

	c.JSON(http.StatusOK, sku)
}

func UpdateSku(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sku ID"})
		return
	}

	var sku models.Sku
	err = c.Bind(&sku)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.UpdateSku(c, id, &sku)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetSku(c, id)
	c.JSON(http.StatusOK, updated)
}