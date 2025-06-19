package controllers

import (
	"net/http"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTenants(c *gin.Context) {
	Tenants, err := models.GetTenants(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Tenants)
}

func GetTenantByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID"})
		return
	}

	Tenant, err := models.GetTenant(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Tenant)
}

func CreateTenant(c *gin.Context) {
	var Tenant models.Tenant

	err := c.Bind(&Tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.CreateTenant(c, &Tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Tenant)
}

func DeleteTenant(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID"})
		return
	}

	Tenant, err := models.DeleteTenant(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, Tenant)
}

func UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tenant ID"})
		return
	}

	var Tenant models.Tenant
	err = c.Bind(&Tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = models.UpdateTenant(c, id, &Tenant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, _ := models.GetTenant(c, id)
	c.JSON(http.StatusOK, updated)
}