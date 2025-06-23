package controllers

import (
	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
)

// GetTenants godoc
// @Summary Get all tenants
// @Tags Tenants
// @Produce json
// @Success 200 {array} models.Tenant
// @Router /tenants [get]
func GetTenants(c *gin.Context) {
	Tenants, err := models.GetTenants(c)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), Tenants)
}

// GetTenantByID godoc
// @Summary Get tenant by ID
// @Tags Tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} models.Tenant
// @Failure 404 {object} map[string]string
// @Router /tenants/{id} [get]
func GetTenantByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Tenant ID")})
		return
	}

	Tenant, err := models.GetTenant(c, id)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusOK), Tenant)
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Tags Tenants
// @Accept json
// @Produce json
// @Param tenant body models.Tenant true "Tenant to create"
// @Success 201 {object} models.Tenant
// @Router /tenants [post]
func CreateTenant(c *gin.Context) {
	var Tenant models.Tenant

	err := c.Bind(&Tenant)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	err = models.CreateTenant(c, &Tenant)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(http.StatusCreated), Tenant)
}

// DeleteTenant godoc
// @Summary Delete tenant by ID
// @Tags Tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} models.Tenant
// @Router /tenants/{id} [delete]
func DeleteTenant(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Tenant ID")})
		return
	}

	Tenant, err := models.DeleteTenant(c, id)
	if err != nil {
		c.JSON(int(http.StatusNotFound), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Tenant not found")})
		return
	}

	c.JSON(int(http.StatusOK), Tenant)
}

// UpdateTenant godoc
// @Summary Update tenant by ID
// @Tags Tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body models.Tenant true "Tenant data"
// @Success 200 {object} models.Tenant
// @Router /tenants/{id} [put]
func UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid Tenant ID")})
		return
	}

	var Tenant models.Tenant
	err = c.Bind(&Tenant)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	err = models.UpdateTenant(c, id, &Tenant)
	if err != nil {
		c.JSON(int(http.StatusInternalServerError), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	updated, _ := models.GetTenant(c, id)
	c.JSON(int(http.StatusOK), updated)
}