package controllers

import (
	"context"
	"errors"
	"testing"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// GetSkus

type mockSkuFetcher struct {
	GetFilteredSkusFunc func(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error)
}

func (m *mockSkuFetcher) GetFilteredSkus(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error) {
	return m.GetFilteredSkusFunc(ctx, tenantID, sellerID, skuCodes)
}

func TestGetSkusLogic(t *testing.T) {
	tenantID := uuid.New()
	sellerID := uuid.New()

	tests := []struct {
		name           string
		tenantIDStr    string
		sellerIDStr    string
		skuCodes       []string
		mockFunc       func(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "missing tenant",
			tenantIDStr:    "",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:           "invalid tenant uuid",
			tenantIDStr:    "bad-uuid",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:           "invalid seller uuid",
			tenantIDStr:    tenantID.String(),
			sellerIDStr:    "bad-seller-id",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "fetch error",
			tenantIDStr: tenantID.String(),
			sellerIDStr: sellerID.String(),
			skuCodes:    []string{"SKU1", "SKU2"},
			mockFunc: func(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:        "success",
			tenantIDStr: tenantID.String(),
			sellerIDStr: sellerID.String(),
			skuCodes:    []string{"SKU1", "SKU2"},
			mockFunc: func(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]models.Sku, error) {
				return []models.Sku{{SkuCode: "SKU1"}, {SkuCode: "SKU2"}}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSkuFetcher{GetFilteredSkusFunc: tt.mockFunc}
			result, status, err := getSkusLogic(mock, tt.tenantIDStr, tt.sellerIDStr, tt.skuCodes)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// GetSkuByID

type mockSkuGetter struct {
	GetSkuFunc func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
}

func (m *mockSkuGetter) GetSku(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
	return m.GetSkuFunc(ctx, id)
}

func TestGetSkuByIDLogic(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid uuid",
			idStr:          "bad-uuid",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "sku not found",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
				return nil, errors.New("sku not found")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
				return &models.Sku{
					ID:       id,
					Name:     "Test SKU",
					SkuCode:  "SKU123",
					SellerID: uuid.New(),
					TenantID: uuid.New(),
				}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSkuGetter{GetSkuFunc: tt.mockFunc}
			result, status, err := getSkuByIDLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Test SKU", result.Name)
			}
		})
	}
}

// CreateSku

type mockSkuCreator struct {
	CreateSkuFunc func(ctx context.Context, sku *models.Sku) error
}

func (m *mockSkuCreator) CreateSku(ctx context.Context, sku *models.Sku) error {
	return m.CreateSkuFunc(ctx, sku)
}

func TestCreateSkuLogic(t *testing.T) {
	validTenantID := uuid.New()
	sku := &models.Sku{Name: "Test", SkuCode: "SKU123", SellerID: uuid.New()}

	tests := []struct {
		name           string
		tenantIDStr    string
		inputSku       *models.Sku
		mockFunc       func(ctx context.Context, sku *models.Sku) error
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "missing tenant id",
			tenantIDStr:    "",
			inputSku:       sku,
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:           "invalid tenant id",
			tenantIDStr:    "invalid-uuid",
			inputSku:       sku,
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "tenant not found",
			tenantIDStr: validTenantID.String(),
			inputSku:    sku,
			mockFunc: func(ctx context.Context, sku *models.Sku) error {
				return gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "creation error",
			tenantIDStr: validTenantID.String(),
			inputSku:    sku,
			mockFunc: func(ctx context.Context, sku *models.Sku) error {
				return errors.New("something went wrong")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:        "success",
			tenantIDStr: validTenantID.String(),
			inputSku:    sku,
			mockFunc: func(ctx context.Context, sku *models.Sku) error {
				return nil
			},
			expectedStatus: int(http.StatusCreated),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSkuCreator{CreateSkuFunc: tt.mockFunc}
			status, err := createSkuLogic(mock, tt.tenantIDStr, tt.inputSku)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, validTenantID, tt.inputSku.TenantID)
			}
		})
	}
}

// DeleteSku

type mockSkuDeleter struct {
	DeleteSkuFunc func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
}

func (m *mockSkuDeleter) DeleteSku(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
	return m.DeleteSkuFunc(ctx, id)
}

func TestDeleteSkuLogic(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid uuid",
			idStr:          "bad-uuid",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "sku not found",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: int(http.StatusNotFound),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
				return &models.Sku{
					ID:       id,
					Name:     "Test SKU",
					SkuCode:  "SKU-001",
					TenantID: uuid.New(),
					SellerID: uuid.New(),
				}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSkuDeleter{DeleteSkuFunc: tt.mockFunc}
			sku, status, err := deleteSkuLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Test SKU", sku.Name)
			}
		})
	}
}

// UpdateSku

type mockSkuUpdater struct {
	UpdateSkuFunc func(ctx context.Context, id uuid.UUID, sku *models.Sku) error
	GetSkuFunc    func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
	GetTenantFunc func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func (m *mockSkuUpdater) UpdateSku(ctx context.Context, id uuid.UUID, sku *models.Sku) error {
	return m.UpdateSkuFunc(ctx, id, sku)
}
func (m *mockSkuUpdater) GetSku(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
	return m.GetSkuFunc(ctx, id)
}
func (m *mockSkuUpdater) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return m.GetTenantFunc(ctx, id)
}

func TestUpdateSkuLogic(t *testing.T) {
	id := uuid.New()
	tenantID := uuid.New()
	sku := &models.Sku{Name: "Updated", TenantID: tenantID}

	tests := []struct {
		name           string
		idStr          string
		sku            *models.Sku
		mockUpdate     func(ctx context.Context, id uuid.UUID, sku *models.Sku) error
		mockGetSku     func(ctx context.Context, id uuid.UUID) (*models.Sku, error)
		mockGetTenant  func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid id",
			idStr:          "bad-id",
			sku:            sku,
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "tenant not found",
			idStr: id.String(),
			sku:   sku,
			mockGetTenant: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "update fails",
			idStr: id.String(),
			sku:   sku,
			mockGetTenant: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{ID: tenantID, Name: "Test"}, nil
			},
			mockUpdate: func(ctx context.Context, id uuid.UUID, sku *models.Sku) error {
				return errors.New("update error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: id.String(),
			sku:   sku,
			mockGetTenant: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{ID: tenantID, Name: "Test"}, nil
			},
			mockUpdate: func(ctx context.Context, id uuid.UUID, sku *models.Sku) error {
				return nil
			},
			mockGetSku: func(ctx context.Context, id uuid.UUID) (*models.Sku, error) {
				return &models.Sku{ID: id, Name: "Updated", SkuCode: "S-001"}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSkuUpdater{
				UpdateSkuFunc: tt.mockUpdate,
				GetSkuFunc:    tt.mockGetSku,
				GetTenantFunc: tt.mockGetTenant,
			}
			res, status, err := updateSkuLogic(mock, tt.idStr, tt.sku)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Updated", res.Name)
			}
		})
	}
}
