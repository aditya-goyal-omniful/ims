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

// GetSellers godoc
// @Summary Get all sellers
// @Tags Sellers
// @Produce json
// @Success 200 {array} models.Seller
// @Router /sellers [get]
func GetSellers(c *gin.Context) {
	sellers, err := models.GetSellers(c)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), sellers)
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
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Seller ID")})
		return
	}

	Seller, err := models.GetSeller(c, id)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), Seller)
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
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	if err := models.CreateSeller(c, &seller); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
			return
		}
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to create seller")})
		return
	}

	c.JSON(int(http.StatusCreated), seller)
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
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Seller ID")})
		return
	}

	seller, err := models.DeleteSeller(c, id)
	if err != nil {
		c.JSON(int(http.StatusNotFound), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Seller not found")})
		return
	}

	c.JSON(int(http.StatusOK), seller)
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
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Seller ID")})
		return
	}

	var seller models.Seller
	err = c.Bind(&seller)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}
	
	if seller.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, seller.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
				return
			}
			c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to validate tenant")})
			return
		}
	}

	err = models.UpdateSeller(c, id, &seller)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	updated, _ := models.GetSeller(c, id)
	c.JSON(int(http.StatusOK), updated)
}