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

// GetSellers

type SellerFetcher interface {
	GetSellers(ctx context.Context) ([]models.Seller, error)
}

func getSellersLogic(service SellerFetcher) ([]models.Seller, int, error) {
	sellers, err := service.GetSellers(context.Background())
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}
	return sellers, int(http.StatusOK), nil
}

// GetSellers godoc
// @Summary Get all sellers
// @Tags Sellers
// @Produce json
// @Success 200 {array} models.Seller
// @Router /sellers [get]
func GetSellers(c *gin.Context) {
	sellers, status, err := getSellersLogic(models.SellerModel{})
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}
	c.JSON(status, sellers)
}

// GetSellerByID

type SellerFetcherByID interface {
	GetSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func getSellerByIDLogic(service SellerFetcherByID, idStr string) (*models.Seller, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid seller id")
	}

	seller, err := service.GetSeller(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	return seller, int(http.StatusOK), nil
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

	seller, status, err := getSellerByIDLogic(models.SellerModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, seller)
}

// CreateSeller

type SellerCreator interface {
	CreateSeller(ctx context.Context, seller *models.Seller) error
}

func createSellerLogic(service SellerCreator, seller *models.Seller) (*models.Seller, int, error) {
	if err := service.CreateSeller(context.Background(), seller); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, int(http.StatusBadRequest), errors.New("tenant not found")
		}
		return nil, int(http.StatusInternalServerError), errors.New("failed to create seller")
	}
	return seller, int(http.StatusCreated), nil
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

	if err := c.Bind(&seller); err != nil {
		c.JSON(int(http.StatusBadRequest), gin.H{i18n.Translate(c, "error"): i18n.Translate(c, "Invalid request body")})
		return
	}

	createdSeller, status, err := createSellerLogic(models.SellerModel{}, &seller)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, createdSeller)
}

// DeleteSeller

type SellerDeleter interface {
	DeleteSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func deleteSellerLogic(service SellerDeleter, idStr string) (*models.Seller, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid seller id")
	}

	seller, err := service.DeleteSeller(context.Background(), id)
	if err != nil {
		return nil, int(http.StatusNotFound), errors.New("seller not found")
	}

	return seller, int(http.StatusOK), nil
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

	seller, status, err := deleteSellerLogic(models.SellerModel{}, idStr)
	if err != nil {
		c.JSON(status, gin.H{i18n.Translate(c, "error"): i18n.Translate(c, err.Error())})
		return
	}

	c.JSON(status, seller)
}

// UpdateSeller

type SellerUpdater interface {
	UpdateSeller(ctx context.Context, id uuid.UUID, seller *models.Seller) error
	GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func updateSellerLogic(service SellerUpdater, idStr string, seller *models.Seller) (*models.Seller, int, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, int(http.StatusBadRequest), errors.New("invalid seller id")
	}

	// Validate Tenant if set
	if seller.TenantID != uuid.Nil {
		if _, err := service.GetTenant(context.Background(), seller.TenantID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, int(http.StatusBadRequest), errors.New("tenant not found")
			}
			return nil, int(http.StatusInternalServerError), errors.New("failed to validate tenant")
		}
	}

	// Update Seller
	if err := service.UpdateSeller(context.Background(), id, seller); err != nil {
		return nil, int(http.StatusInternalServerError), err
	}

	// Return updated
	updated, _ := service.GetSeller(context.Background(), id)
	return updated, int(http.StatusOK), nil
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