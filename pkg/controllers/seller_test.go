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

// GetSellers

type mockSellerFetcher struct {
	GetSellersFunc func(ctx context.Context) ([]models.Seller, error)
}

func (m *mockSellerFetcher) GetSellers(ctx context.Context) ([]models.Seller, error) {
	return m.GetSellersFunc(ctx)
}

func TestGetSellersLogic(t *testing.T) {
	mockSellers := []models.Seller{
		{ID: uuid.New(), Name: "Seller A"},
		{ID: uuid.New(), Name: "Seller B"},
	}

	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]models.Seller, error)
		expectedCount  int
		expectedStatus int
		expectErr      bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]models.Seller, error) {
				return mockSellers, nil
			},
			expectedCount:  len(mockSellers),
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
		{
			name: "db error",
			mockFunc: func(ctx context.Context) ([]models.Seller, error) {
				return nil, errors.New("DB failure")
			},
			expectedCount:  0,
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSellerFetcher{GetSellersFunc: tt.mockFunc}
			sellers, status, err := getSellersLogic(mock)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedCount, len(sellers))
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// GetSellerByID

type mockSellerByIDFetcher struct {
	GetSellerFunc func(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func (m *mockSellerByIDFetcher) GetSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	return m.GetSellerFunc(ctx, id)
}

func TestGetSellerByIDLogic(t *testing.T) {
	id := uuid.New()
	mockSeller := &models.Seller{ID: id, Name: "Test Seller"}

	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Seller, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid uuid",
			idStr: "invalid-uuid",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return nil, nil
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "seller not found (db error)",
			idStr: id.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "successful fetch",
			idStr: id.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return mockSeller, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSellerByIDFetcher{GetSellerFunc: tt.mockFunc}
			seller, status, err := getSellerByIDLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, seller)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mockSeller.ID, seller.ID)
			}
		})
	}
}

// CreateSeller

type mockSellerCreator struct {
	CreateSellerFunc func(ctx context.Context, seller *models.Seller) error
}

func (m *mockSellerCreator) CreateSeller(ctx context.Context, seller *models.Seller) error {
	return m.CreateSellerFunc(ctx, seller)
}

func TestCreateSellerLogic(t *testing.T) {
	newSeller := &models.Seller{ID: uuid.New(), Name: "New Seller"}

	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context, seller *models.Seller) error
		expectedStatus int
		expectErr      bool
	}{
		{
			name: "tenant not found",
			mockFunc: func(ctx context.Context, seller *models.Seller) error {
				return gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name: "create fails",
			mockFunc: func(ctx context.Context, seller *models.Seller) error {
				return errors.New("db error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name: "create succeeds",
			mockFunc: func(ctx context.Context, seller *models.Seller) error {
				return nil
			},
			expectedStatus: int(http.StatusCreated),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSellerCreator{CreateSellerFunc: tt.mockFunc}
			result, status, err := createSellerLogic(mock, newSeller)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, newSeller.ID, result.ID)
			}
		})
	}
}

// DeleteSeller

type mockSellerDeleter struct {
	DeleteSellerFunc func(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func (m *mockSellerDeleter) DeleteSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	return m.DeleteSellerFunc(ctx, id)
}

func TestDeleteSellerLogic(t *testing.T) {
	sellerID := uuid.New()
	mockSeller := &models.Seller{ID: sellerID, Name: "SellerX"}

	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Seller, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid uuid",
			idStr: "invalid-id",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return nil, nil
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "not found",
			idStr: sellerID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return nil, errors.New("not found")
			},
			expectedStatus: int(http.StatusNotFound),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: sellerID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
				return mockSeller, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSellerDeleter{DeleteSellerFunc: tt.mockFunc}
			seller, status, err := deleteSellerLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, seller)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mockSeller.ID, seller.ID)
			}
		})
	}
}

// UpdateSeller

type mockSellerUpdater struct {
	UpdateSellerFunc func(ctx context.Context, id uuid.UUID, seller *models.Seller) error
	GetTenantFunc    func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetSellerFunc    func(ctx context.Context, id uuid.UUID) (*models.Seller, error)
}

func (m *mockSellerUpdater) UpdateSeller(ctx context.Context, id uuid.UUID, seller *models.Seller) error {
	return m.UpdateSellerFunc(ctx, id, seller)
}

func (m *mockSellerUpdater) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return m.GetTenantFunc(ctx, id)
}

func (m *mockSellerUpdater) GetSeller(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	return m.GetSellerFunc(ctx, id)
}

func TestUpdateSellerLogic(t *testing.T) {
	validID := uuid.New()
	seller := &models.Seller{Name: "SellerX", TenantID: validID}
	updatedSeller := &models.Seller{ID: validID, Name: "Updated Seller"}

	tests := []struct {
		name           string
		idStr          string
		mockUpdater    *mockSellerUpdater
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid UUID",
			idStr: "bad-id",
			mockUpdater: &mockSellerUpdater{},
			expectedStatus: int(http.StatusBadRequest),
			expectErr: true,
		},
		{
			name:  "tenant not found",
			idStr: validID.String(),
			mockUpdater: &mockSellerUpdater{
				GetTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "update fails",
			idStr: validID.String(),
			mockUpdater: &mockSellerUpdater{
				GetTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
					return &models.Tenant{ID: id}, nil
				},
				UpdateSellerFunc: func(ctx context.Context, id uuid.UUID, seller *models.Seller) error {
					return errors.New("db failure")
				},
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockUpdater: &mockSellerUpdater{
				GetTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
					return &models.Tenant{ID: id}, nil
				},
				UpdateSellerFunc: func(ctx context.Context, id uuid.UUID, seller *models.Seller) error {
					return nil
				},
				GetSellerFunc: func(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
					return updatedSeller, nil
				},
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, status, err := updateSellerLogic(tt.mockUpdater, tt.idStr, seller)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, updatedSeller.Name, resp.Name)
			}
		})
	}
}
