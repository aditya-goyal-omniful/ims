package controllers

import (
	"context"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
)

// GetHubs

type TenantFetcher interface {
	GetAllTenants(ctx context.Context) ([]models.Tenant, error)
}

func getTenantsLogic(service TenantFetcher) ([]models.Tenant, int) {
	tenants, err := service.GetAllTenants(context.Background())
	if err != nil {
		return nil, int(http.StatusInternalServerError)
	}
	return tenants, int(http.StatusOK)
}

// GetTenants godoc
// @Summary Get all tenants
// @Tags Tenants
// @Produce json
// @Success 200 {array} models.Tenant
// @Router /tenants [get]
func GetTenants(c *gin.Context) {
	tenants, status := getTenantsLogic(models.TenantModel{})

	if status != int(http.StatusOK) {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to fetch tenants")})
		return
	}
	c.JSON(status, tenants)
}

// GetTenantByID

type TenantService interface {
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func getTenantByIDLogic(service TenantService, idStr string) (*models.Tenant, int) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest)
	}

	tenant, err := service.GetTenant(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError)
	}

	return tenant, int(http.StatusOK)
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

	tenant, status := getTenantByIDLogic(models.TenantModel{}, idStr)

	if status != int(http.StatusOK) {
		msg := "Error fetching tenant"
		if status == int(http.StatusBadRequest) {
			msg = "Invalid tenant ID"
		}
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, msg)})
		return
	}

	c.JSON(int(http.StatusOK), tenant)
}

// CreateTenant

type TenantCreator interface {
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
}

func createTenantLogic(service TenantCreator, tenant *models.Tenant) (int, error) {
	err := service.CreateTenant(context.Background(), tenant)
	if err != nil {
		return int(http.StatusInternalServerError), err
	}
	return int(http.StatusCreated), nil
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
	var tenant models.Tenant

	if err := c.Bind(&tenant); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	status, err := createTenantLogic(models.TenantModel{}, &tenant)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, tenant)
}

// DeleteTenant

type TenantDeleter interface {
	DeleteTenant(ctx context.Context, id uuid.UUID) (models.Tenant, error)
}

func deleteTenantLogic(service TenantDeleter, idStr string) (models.Tenant, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Tenant{}, int(http.StatusBadRequest), err
	}

	tenant, err := service.DeleteTenant(context.Background(), id)
	if err != nil {
		return models.Tenant{}, int(http.StatusNotFound), err
	}

	return tenant, int(http.StatusOK), nil
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

	tenant, status, err := deleteTenantLogic(models.TenantModel{}, idStr)
	if err != nil {
		msg := "Tenant not found"
		if status == int(http.StatusBadRequest) {
			msg = "Invalid Tenant ID"
		}
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, msg)})
		return
	}

	c.JSON(int(http.StatusOK), tenant)
}

// UpdateTenant

type TenantUpdater interface {
	UpdateTenant(ctx context.Context, id uuid.UUID, updated *models.Tenant) error
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func updateTenantLogic(service TenantUpdater, idStr string, updated *models.Tenant) (*models.Tenant, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), err
	}

	err = service.UpdateTenant(context.Background(), id, updated)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	tenant, err := service.GetTenant(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return tenant, int(http.StatusOK), nil
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

	var tenant models.Tenant
	if err := c.Bind(&tenant); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	updatedTenant, status, err := updateTenantLogic(models.TenantModel{}, idStr, &tenant)
	if err != nil {
		msg := "Error updating tenant"
		if status ==int(http.StatusBadRequest) {
			msg = "Invalid Tenant ID"
		}
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, msg)})
		return
	}

	c.JSON(status, updatedTenant)
}
