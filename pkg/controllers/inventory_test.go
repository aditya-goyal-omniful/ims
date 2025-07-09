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

// GetInventories

type mockInventoryFetcher struct {
	GetInventoriesFunc func(ctx context.Context) ([]models.Inventory, error)
}

func (m *mockInventoryFetcher) GetInventories(ctx context.Context) ([]models.Inventory, error) {
	return m.GetInventoriesFunc(ctx)
}

func TestGetInventoriesLogic(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]models.Inventory, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name: "fetch error",
			mockFunc: func(ctx context.Context) ([]models.Inventory, error) {
				return nil, errors.New("DB error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]models.Inventory, error) {
				return []models.Inventory{
					{SkuID: uuid.New(), HubID: uuid.New(), Quantity: 50},
				}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryFetcher{GetInventoriesFunc: tt.mockFunc}
			result, status, err := getInventoriesLogic(mock)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
			}
		})
	}
}

// GetInventoryByID

type mockInventoryByIDFetcher struct {
	GetInventoryFunc func(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
}

func (m *mockInventoryByIDFetcher) GetInventory(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	return m.GetInventoryFunc(ctx, id)
}

func TestGetInventoryByIDLogic(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid uuid",
			idStr:          "not-a-uuid",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "db error",
			idStr: id.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: id.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
				return &models.Inventory{ID: id, Quantity: 100}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryByIDFetcher{GetInventoryFunc: tt.mockFunc}
			result, status, err := getInventoryByIDLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 100, result.Quantity)
			}
		})
	}
}

// CreateInventory

type mockInventoryCreator struct {
	CreateInventoryFunc func(ctx context.Context, inv *models.Inventory) error
}

func (m *mockInventoryCreator) CreateInventory(ctx context.Context, inv *models.Inventory) error {
	return m.CreateInventoryFunc(ctx, inv)
}

func TestCreateInventoryLogic(t *testing.T) {
	validTenantID := uuid.New().String()

	tests := []struct {
		name           string
		tenantID       string
		inv            *models.Inventory
		mockFunc       func(ctx context.Context, inv *models.Inventory) error
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid tenant id",
			tenantID:       "not-a-uuid",
			inv:            &models.Inventory{},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:     "tenant not found",
			tenantID: validTenantID,
			inv:      &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:     "db error",
			tenantID: validTenantID,
			inv:      &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return errors.New("db error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:     "success",
			tenantID: validTenantID,
			inv:      &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return nil
			},
			expectedStatus: int(http.StatusCreated),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryCreator{CreateInventoryFunc: tt.mockFunc}
			status, err := createInventoryLogic(mock, tt.tenantID, tt.inv)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// DeleteInventory

type mockInventoryUpdater struct {
	UpdateInventoryFunc func(ctx context.Context, id uuid.UUID, inv *models.Inventory) error
	GetInventoryFunc    func(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
}

func (m *mockInventoryUpdater) UpdateInventory(ctx context.Context, id uuid.UUID, inv *models.Inventory) error {
	return m.UpdateInventoryFunc(ctx, id, inv)
}
func (m *mockInventoryUpdater) GetInventory(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	return m.GetInventoryFunc(ctx, id)
}

type mockTenantValidator struct {
	GetTenantFunc func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func (m *mockTenantValidator) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return m.GetTenantFunc(ctx, id)
}

func TestUpdateInventoryLogic(t *testing.T) {
	id := uuid.New()
	validTenantID := uuid.New()

	tests := []struct {
		name           string
		idStr          string
		inv            *models.Inventory
		updateFunc     func(ctx context.Context, id uuid.UUID, inv *models.Inventory) error
		getFunc        func(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
		tenantFunc     func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid id",
			idStr:          "invalid-uuid",
			inv:            &models.Inventory{},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "tenant not found",
			idStr: id.String(),
			inv:   &models.Inventory{TenantID: validTenantID},
			tenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:  "update error",
			idStr: id.String(),
			inv:   &models.Inventory{},
			tenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{}, nil
			},
			updateFunc: func(ctx context.Context, id uuid.UUID, inv *models.Inventory) error {
				return errors.New("db update error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: id.String(),
			inv:   &models.Inventory{Quantity: 10},
			tenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{}, nil
			},
			updateFunc: func(ctx context.Context, id uuid.UUID, inv *models.Inventory) error {
				return nil
			},
			getFunc: func(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
				return &models.Inventory{ID: id, Quantity: 10}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := &mockInventoryUpdater{
				UpdateInventoryFunc: tt.updateFunc,
				GetInventoryFunc:    tt.getFunc,
			}
			tenantValidator := &mockTenantValidator{
				GetTenantFunc: tt.tenantFunc,
			}

			result, status, err := updateInventoryLogic(updater, tenantValidator, tt.idStr, tt.inv)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 10, result.Quantity)
			}
		})
	}
}

// UpsertInventory

type mockInventoryUpserter struct {
	UpsertInventoryFunc func(ctx context.Context, inv *models.Inventory) error
}

func (m *mockInventoryUpserter) UpsertInventory(ctx context.Context, inv *models.Inventory) error {
	return m.UpsertInventoryFunc(ctx, inv)
}

func TestUpsertInventoryLogic(t *testing.T) {
	validTenant := uuid.New()
	tests := []struct {
		name           string
		tenantIDStr    string
		input          *models.Inventory
		mockFunc       func(ctx context.Context, inv *models.Inventory) error
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid tenant ID",
			tenantIDStr:    "not-a-uuid",
			input:          &models.Inventory{},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "tenant not found",
			tenantIDStr: validTenant.String(),
			input:       &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return gorm.ErrRecordNotFound
			},
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "db error",
			tenantIDStr: validTenant.String(),
			input:       &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return errors.New("DB failure")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:        "success",
			tenantIDStr: validTenant.String(),
			input:       &models.Inventory{},
			mockFunc: func(ctx context.Context, inv *models.Inventory) error {
				return nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryUpserter{
				UpsertInventoryFunc: tt.mockFunc,
			}
			status, err := upsertInventoryLogic(mock, tt.tenantIDStr, tt.input)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ViewInventoryWithDefault

type mockInventoryViewer struct {
	GetInventoryWithDefaultsFunc func(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error)
}

func (m *mockInventoryViewer) GetInventoryWithDefaults(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error) {
	return m.GetInventoryWithDefaultsFunc(ctx, tenantID, hubID)
}

func TestViewInventoryWithDefaultsLogic(t *testing.T) {
	validTenantID := uuid.New()
	validHubID := uuid.New()

	tests := []struct {
		name           string
		tenantIDStr    string
		hubIDStr       string
		mockFunc       func(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid tenant ID",
			tenantIDStr:    "bad-uuid",
			hubIDStr:       validHubID.String(),
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:           "invalid hub ID",
			tenantIDStr:    validTenantID.String(),
			hubIDStr:       "bad-uuid",
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:        "DB error",
			tenantIDStr: validTenantID.String(),
			hubIDStr:    validHubID.String(),
			mockFunc: func(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error) {
				return nil, errors.New("DB error")
			},
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:        "success",
			tenantIDStr: validTenantID.String(),
			hubIDStr:    validHubID.String(),
			mockFunc: func(ctx context.Context, tenantID, hubID uuid.UUID) ([]models.InventoryView, error) {
				return []models.InventoryView{{SkuCode: "sku-123", Quantity: 0}}, nil
			},
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryViewer{
				GetInventoryWithDefaultsFunc: tt.mockFunc,
			}
			result, status, err := viewInventoryWithDefaultsLogic(mock, tt.tenantIDStr, tt.hubIDStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

// CheckAndUpdateInventory

type mockInventoryChecker struct {
	GetInventoryBySkuHubFunc    func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error)
	UpdateInventoryQuantityFunc func(ctx context.Context, invID uuid.UUID, newQty int) error
}

func (m *mockInventoryChecker) GetInventoryBySkuHub(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
	return m.GetInventoryBySkuHubFunc(ctx, skuID, hubID)
}

func (m *mockInventoryChecker) UpdateInventoryQuantity(ctx context.Context, invID uuid.UUID, newQty int) error {
	return m.UpdateInventoryQuantityFunc(ctx, invID, newQty)
}

func TestCheckAndUpdateInventoryLogic(t *testing.T) {
	skuID := uuid.New()
	hubID := uuid.New()
	invID := uuid.New()

	tests := []struct {
		name           string
		req            CheckInventoryRequest
		mockFetch      func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error)
		mockUpdate     func(ctx context.Context, invID uuid.UUID, newQty int) error
		expectedAvail  bool
		expectedStatus int
		expectErr      bool
	}{
		{
			name: "inventory not found",
			req: CheckInventoryRequest{
				SKUID:    skuID,
				HubID:    hubID,
				Quantity: 5,
			},
			mockFetch: func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
				return nil, gorm.ErrRecordNotFound
			},
			mockUpdate:     nil,
			expectedAvail:  false,
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
		{
			name: "inventory fetch error",
			req:  CheckInventoryRequest{SKUID: skuID, HubID: hubID, Quantity: 5},
			mockFetch: func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
				return nil, errors.New("DB error")
			},
			mockUpdate:     nil,
			expectedAvail:  false,
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name: "insufficient inventory",
			req:  CheckInventoryRequest{SKUID: skuID, HubID: hubID, Quantity: 10},
			mockFetch: func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
				return &models.Inventory{ID: invID, Quantity: 5}, nil
			},
			mockUpdate:     nil,
			expectedAvail:  false,
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
		{
			name: "update error",
			req:  CheckInventoryRequest{SKUID: skuID, HubID: hubID, Quantity: 5},
			mockFetch: func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
				return &models.Inventory{ID: invID, Quantity: 10}, nil
			},
			mockUpdate: func(ctx context.Context, invID uuid.UUID, newQty int) error {
				return errors.New("update fail")
			},
			expectedAvail:  false,
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name: "success",
			req:  CheckInventoryRequest{SKUID: skuID, HubID: hubID, Quantity: 3},
			mockFetch: func(ctx context.Context, skuID, hubID uuid.UUID) (*models.Inventory, error) {
				return &models.Inventory{ID: invID, Quantity: 10}, nil
			},
			mockUpdate: func(ctx context.Context, invID uuid.UUID, newQty int) error {
				assert.Equal(t, 7, newQty)
				return nil
			},
			expectedAvail:  true,
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInventoryChecker{
				GetInventoryBySkuHubFunc:    tt.mockFetch,
				UpdateInventoryQuantityFunc: tt.mockUpdate,
			}
			ok, status, err := checkAndUpdateInventoryLogic(mock, tt.req)
			assert.Equal(t, tt.expectedAvail, ok)
			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
