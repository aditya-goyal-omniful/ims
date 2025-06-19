package controllers

import (
	"errors"
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSkus(c *gin.Context) {
	skus, err := models.GetSkus(c)
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

	if err := models.CreateSku(c, &sku); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant or Seller not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sku"})
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

	if sku.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, sku.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate tenant"})
			return
		}
	}

	err = models.UpdateSku(c, id, &sku)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetSku(c, id)
	c.JSON(http.StatusOK, updated)
}