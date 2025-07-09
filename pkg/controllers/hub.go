package controllers

import (
	"context"
	"errors"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"gorm.io/gorm"
)

// GetHubs

type HubFetcher interface {
	GetAllHubs(ctx context.Context) ([]models.Hub, error)
}

func getHubsLogic(service HubFetcher) ([]models.Hub, int) {
	hubs, err := service.GetAllHubs(context.Background())
	if err != nil {
		return nil, int(http.StatusInternalServerError)
	}
	return hubs, int(http.StatusOK)
}

// GetHubs godoc
// @Summary Get all hubs
// @Tags Hubs
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {array} models.Hub
// @Router /hubs [get]
func GetHubs(c *gin.Context) {
	hubs, status := getHubsLogic(models.HubModel{})

	if status != int(http.StatusOK) {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Failed to fetch hubs")})
		return
	}
	c.JSON(status, hubs)
}

// GetHubByID

type HubService interface {
	GetHub(ctx context.Context, id uuid.UUID) (*models.Hub, error)
}

func getHubByIDLogic(service HubService, idStr string) (*models.Hub, int) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest)
	}

	hub, err := service.GetHub(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError)
	}

	return hub, int(http.StatusOK)
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

	hub, status := getHubByIDLogic(models.HubModel{}, idStr)

	if status != int(http.StatusOK) {
		msg := "Error fetching hub"
		if status == int(http.StatusBadRequest) {
			msg = "Invalid hub ID"
		}
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, msg)})
		return
	}

	c.JSON(int(http.StatusOK), hub)
}

// CreateHub

type HubCreator interface {
	CreateHub(ctx context.Context, hub *models.Hub) error
}

func createHubLogic(service HubCreator, tenantIDStr string, hub *models.Hub) (int, error) {
	// Parse tenant ID
	if tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			return int(http.StatusBadRequest), errors.New("invalid tenant_id")
		}
		hub.TenantID = tenantID
	}

	// Create hub
	err := service.CreateHub(context.Background(), hub)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return int(http.StatusBadRequest), errors.New("tenant not found")
		}
		return int(http.StatusInternalServerError), errors.New("failed to create hub")
	}

	return int(http.StatusCreated), nil
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

	if err := c.Bind(&hub); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	status, err := createHubLogic(models.HubModel{}, c.GetHeader("X-Tenant-ID"), &hub)
	if err != nil {
		c.JSON(int(status), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(int(status), hub)
}

// DeleteHub

type HubDeleter interface {
	DeleteHub(ctx context.Context, id uuid.UUID) (models.Hub, error)
}

func deleteHubLogic(service HubDeleter, idStr string) (models.Hub, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return models.Hub{}, int(http.StatusBadRequest), errors.New("invalid hub id")
	}

	hub, err := service.DeleteHub(context.Background(), id)
	if err != nil {
		return models.Hub{}, int(http.StatusNotFound), errors.New("hub not found")
	}

	return hub, int(http.StatusOK), nil
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
	hub, status, err := deleteHubLogic(models.HubModel{}, c.Param("id"))
	if err != nil {
		c.JSON(int(status), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}
	c.JSON(int(status), hub)
}

// UpdateHub

func updateHubLogic(
	ctx context.Context,
	tenantService TenantService,
	updateFunc func(ctx context.Context, id uuid.UUID, hub *models.Hub) error,
	getFunc func(ctx context.Context, id uuid.UUID) (*models.Hub, error),
	idStr string,
	input models.Hub,
) (*models.Hub, int) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest)
	}

	if input.TenantID != uuid.Nil && tenantService != nil {
		_, err := tenantService.GetTenant(ctx, input.TenantID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, int(http.StatusBadRequest)
			}
			return nil, int(http.StatusInternalServerError)
		}
	}

	if err := updateFunc(ctx, id, &input); err != nil {
		return nil, int(http.StatusInternalServerError)
	}

	updated, err := getFunc(ctx, id)
	if err != nil {
		return nil, int(http.StatusInternalServerError)
	}

	return updated, int(http.StatusOK)
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

	var hub models.Hub
	err := c.Bind(&hub)
	if err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	result, status := updateHubLogic(
		c,
		models.TenantModel{},
		models.UpdateHub,
		models.GetHub,
		idStr,
		hub,
	)

		if status != int(http.StatusOK) {
		msg := "Failed to update hub"
		if status == int(http.StatusBadRequest) {
			msg = "Invalid hub ID or tenant not found"
		}
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, msg)})
		return
	}

	c.JSON(int(http.StatusOK), result)
}