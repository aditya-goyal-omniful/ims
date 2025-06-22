package controllers

import (
	"errors"
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetSellers godoc
// @Summary Get all sellers
// @Tags Sellers
// @Produce json
// @Success 200 {array} models.Seller
// @Router /sellers [get]
func GetSellers(c *gin.Context) {
	sellers, err := models.GetSellers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sellers)
}

// GetSellerByID godoc
// @Summary Get seller by ID
// @Tags Sellers
// @Produce json
// @Param id path string true "Seller ID"
// @Success 200 {object} models.Seller
// @Router /sellers/{id} [get]
func GetSellerByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Seller ID"})
		return
	}

	Seller, err := models.GetSeller(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Seller)
}

// CreateSeller godoc
// @Summary Create a new seller
// @Tags Sellers
// @Accept json
// @Produce json
// @Param seller body models.Seller true "Seller to create"
// @Success 201 {object} models.Seller
// @Router /sellers [post]
func CreateSeller(c *gin.Context) {
	var seller models.Seller

	err := c.Bind(&seller)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := models.CreateSeller(c, &seller); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seller"})
		return
	}

	c.JSON(http.StatusCreated, seller)
}

// DeleteSeller godoc
// @Summary Delete seller by ID
// @Tags Sellers
// @Produce json
// @Param id path string true "Seller ID"
// @Success 200 {object} models.Seller
// @Router /sellers/{id} [delete]
func DeleteSeller(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Seller ID"})
		return
	}

	seller, err := models.DeleteSeller(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
		return
	}

	c.JSON(http.StatusOK, seller)
}

// UpdateSeller godoc
// @Summary Update seller by ID
// @Tags Sellers
// @Accept json
// @Produce json
// @Param id path string true "Seller ID"
// @Param seller body models.Seller true "Updated seller"
// @Success 200 {object} models.Seller
// @Router /sellers/{id} [put]
func UpdateSeller(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Seller ID"})
		return
	}

	var seller models.Seller
	err = c.Bind(&seller)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	if seller.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, seller.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate tenant"})
			return
		}
	}

	err = models.UpdateSeller(c, id, &seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetSeller(c, id)
	c.JSON(http.StatusOK, updated)
}