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

// GetHubs godoc
// @Summary Get all hubs
// @Tags Hubs
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {array} models.Hub
// @Router /hubs [get]
func GetHubs(c *gin.Context) {
	hubs, err := models.GetHubs(c)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), hubs)
}

// GetHubByID godoc
// @Summary Get hub by ID
// @Tags Hubs
// @Produce json
// @Param id path string true "Hub ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Hub
// @Router /hubs/{id} [get]
func GetHubByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid hub ID")})
		return
	}

	hub, err := models.GetHub(c, id)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), hub)
}

// CreateHub godoc
// @Summary Create a new hub
// @Tags Hubs
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param hub body models.Hub true "Hub to create"
// @Success 201 {object} models.Hub
// @Router /hubs [post]
func CreateHub(c *gin.Context) {
	var hub models.Hub

	err := c.Bind(&hub)
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
		hub.TenantID = tenantID
	}

	if err := models.CreateHub(c, &hub); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
			return
		}
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to create hub")})
		return
	}

	c.JSON(int(http.StatusCreated), hub)
}

// DeleteHub godoc
// @Summary Delete hub by ID
// @Tags Hubs
// @Produce json
// @Param id path string true "Hub ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} models.Hub
// @Router /hubs/{id} [delete]
func DeleteHub(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid hub ID")})
		return
	}

	hub, err := models.DeleteHub(c, id)
	if err != nil {
		c.JSON(int(http.StatusNotFound), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Hub not found")})
		return
	}

	c.JSON(int(http.StatusOK), hub)
}

// UpdateHub godoc
// @Summary Update hub by ID
// @Tags Hubs
// @Accept json
// @Produce json
// @Param id path string true "Hub ID"
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param hub body models.Hub true "Updated hub"
// @Success 200 {object} models.Hub
// @Router /hubs/{id} [put]
func UpdateHub(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid hub ID")})
		return
	}

	var hub models.Hub
	err = c.Bind(&hub)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	if hub.TenantID != uuid.Nil {
		if _, err := models.GetTenant(c, hub.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
				return
			}
			c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to validate tenant")})
			return
		}
	}

	err = models.UpdateHub(c, id, &hub)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	updated, _ := models.GetHub(c, id)
	c.JSON(int(http.StatusOK), updated)
}