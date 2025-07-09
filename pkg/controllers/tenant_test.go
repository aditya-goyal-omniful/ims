package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// GetTenants

type mockTenantFetcher struct {
	GetAllTenantsFunc func(ctx context.Context) ([]models.Tenant, error)
}

func (m *mockTenantFetcher) GetAllTenants(ctx context.Context) ([]models.Tenant, error) {
	return m.GetAllTenantsFunc(ctx)
}

func TestGetTenantsLogic(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func(ctx context.Context) ([]models.Tenant, error)
		expectedStatus int
		expectedLength int
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]models.Tenant, error) {
				return []models.Tenant{
					{ID: uuid.New(), Name: "Tenant One"},
					{ID: uuid.New(), Name: "Tenant Two"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedLength: 2,
		},
		{
			name: "error from service",
			mockFunc: func(ctx context.Context) ([]models.Tenant, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTenantFetcher{GetAllTenantsFunc: tt.mockFunc}
			tenants, status := getTenantsLogic(mock)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Len(t, tenants, tt.expectedLength)
		})
	}
}

// GetTenantById

type mockTenantService struct {
	GetTenantFunc func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func (m *mockTenantService) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return m.GetTenantFunc(ctx, id)
}

func TestGetTenantByIDLogic(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
		expectedStatus int
		expectNil      bool
	}{
		{
			name:  "invalid uuid",
			idStr: "not-a-uuid",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, nil // shouldn't be called
			},
			expectedStatus: http.StatusBadRequest,
			expectNil:      true,
		},
		{
			name:  "service error",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, errors.New("db error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectNil:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{ID: id, Name: "Tenant X"}, nil
			},
			expectedStatus: http.StatusOK,
			expectNil:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTenantService{GetTenantFunc: tt.mockFunc}
			result, status := getTenantByIDLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

// CreateTenant

type mockTenantCreator struct {
	CreateTenantFunc func(ctx context.Context, tenant *models.Tenant) error
}

func (m *mockTenantCreator) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return m.CreateTenantFunc(ctx, tenant)
}

func TestCreateTenantLogic(t *testing.T) {
	tests := []struct {
		name           string
		input          *models.Tenant
		mockFunc       func(ctx context.Context, tenant *models.Tenant) error
		expectedStatus int
		expectErr      bool
	}{
		{
			name: "success",
			input: &models.Tenant{
				Name: "TestTenant",
			},
			mockFunc: func(ctx context.Context, tenant *models.Tenant) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectErr:      false,
		},
		{
			name: "creation failed",
			input: &models.Tenant{
				Name: "TestTenant",
			},
			mockFunc: func(ctx context.Context, tenant *models.Tenant) error {
				return errors.New("db failure")
			},
			expectedStatus: http.StatusInternalServerError,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTenantCreator{CreateTenantFunc: tt.mockFunc}
			status, err := createTenantLogic(mock, tt.input)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// DeleteTenant

type mockTenantDeleter struct {
	DeleteTenantFunc func(ctx context.Context, id uuid.UUID) (models.Tenant, error)
}

func (m *mockTenantDeleter) DeleteTenant(ctx context.Context, id uuid.UUID) (models.Tenant, error) {
	return m.DeleteTenantFunc(ctx, id)
}

func TestDeleteTenantLogic(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name           string
		idStr          string
		mockFunc       func(ctx context.Context, id uuid.UUID) (models.Tenant, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid uuid",
			idStr: "bad-uuid",
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Tenant, error) {
				return models.Tenant{}, nil // should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectErr:      true,
		},
		{
			name:  "tenant not found",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Tenant, error) {
				return models.Tenant{}, errors.New("not found")
			},
			expectedStatus: http.StatusNotFound,
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Tenant, error) {
				return models.Tenant{ID: id, Name: "Deleted Tenant"}, nil
			},
			expectedStatus: http.StatusOK,
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTenantDeleter{DeleteTenantFunc: tt.mockFunc}
			tenant, status, err := deleteTenantLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.idStr, tenant.ID.String())
			}
		})
	}
}

// UpdateTenant

type mockTenantUpdater struct {
	UpdateTenantFunc func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error
	GetTenantFunc    func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
}

func (m *mockTenantUpdater) UpdateTenant(ctx context.Context, id uuid.UUID, updated *models.Tenant) error {
	return m.UpdateTenantFunc(ctx, id, updated)
}

func (m *mockTenantUpdater) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	return m.GetTenantFunc(ctx, id)
}

func TestUpdateTenantLogic(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name           string
		idStr          string
		input          *models.Tenant
		updateFunc     func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error
		getFunc        func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid uuid",
			idStr: "bad-uuid",
			input: &models.Tenant{Name: "TenantA"},
			updateFunc: func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error {
				return nil // shouldn't be called
			},
			getFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, nil // shouldn't be called
			},
			expectedStatus: http.StatusBadRequest,
			expectErr:      true,
		},
		{
			name:  "update failed",
			idStr: validID.String(),
			input: &models.Tenant{Name: "TenantB"},
			updateFunc: func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error {
				return errors.New("update error")
			},
			getFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectErr:      true,
		},
		{
			name:  "get after update failed",
			idStr: validID.String(),
			input: &models.Tenant{Name: "TenantC"},
			updateFunc: func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error {
				return nil
			},
			getFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, errors.New("get error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			input: &models.Tenant{Name: "TenantD"},
			updateFunc: func(ctx context.Context, id uuid.UUID, updated *models.Tenant) error {
				return nil
			},
			getFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{ID: id, Name: "TenantD"}, nil
			},
			expectedStatus: http.StatusOK,
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTenantUpdater{
				UpdateTenantFunc: tt.updateFunc,
				GetTenantFunc:    tt.getFunc,
			}
			tenant, status, err := updateTenantLogic(mock, tt.idStr, tt.input)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Name, tenant.Name)
			}
		})
	}
}
