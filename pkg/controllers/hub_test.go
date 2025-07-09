package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// GetHubs

type mockHubFetcher struct {
	GetAllHubsFunc func(ctx context.Context) ([]models.Hub, error)
}

func (m *mockHubFetcher) GetAllHubs(ctx context.Context) ([]models.Hub, error) {
	return m.GetAllHubsFunc(ctx)
}

func TestGetHubsLogic(t *testing.T) {
	mockHubs := []models.Hub{
		{ID: uuid.New(), Name: "Hub A"},
		{ID: uuid.New(), Name: "Hub B"},
	}

	tests := []struct {
		name          string
		mockFunc      func(ctx context.Context) ([]models.Hub, error)
		expectedHubs  []models.Hub
		expectedCode  int
		expectFailure bool
	}{
		{
			name: "success - hubs returned",
			mockFunc: func(ctx context.Context) ([]models.Hub, error) {
				return mockHubs, nil
			},
			expectedHubs:  mockHubs,
			expectedCode:  http.StatusOK,
			expectFailure: false,
		},
		{
			name: "failure - internal error",
			mockFunc: func(ctx context.Context) ([]models.Hub, error) {
				return nil, errors.New("DB error")
			},
			expectedHubs:  nil,
			expectedCode:  http.StatusInternalServerError,
			expectFailure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHubFetcher{
				GetAllHubsFunc: tt.mockFunc,
			}

			hubs, status := getHubsLogic(mock)

			assert.Equal(t, tt.expectedCode, status)
			if tt.expectFailure {
				assert.Nil(t, hubs)
			} else {
				assert.Equal(t, tt.expectedHubs, hubs)
			}
		})
	}
}

// GetHubById

type mockHubService struct {
	GetHubFunc func(ctx context.Context, id uuid.UUID) (*models.Hub, error)
}

func (m *mockHubService) GetHub(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
	return m.GetHubFunc(ctx, id)
}

func TestGetHubByIDLogic(t *testing.T) {
	validID := uuid.New()
	expectedHub := &models.Hub{ID: validID, Name: "Main Hub"}

	tests := []struct {
		name         string
		idStr        string
		mockFunc     func(ctx context.Context, id uuid.UUID) (*models.Hub, error)
		expectedHub  *models.Hub
		expectedCode int
	}{
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
				return expectedHub, nil
			},
			expectedHub:  expectedHub,
			expectedCode: http.StatusOK,
		},
		{
			name:  "invalid UUID",
			idStr: "invalid-uuid",
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
				return nil, nil // won't be called
			},
			expectedHub:  nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "DB error",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
				return nil, errors.New("DB failure")
			},
			expectedHub:  nil,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHubService{
				GetHubFunc: tt.mockFunc,
			}

			hub, status := getHubByIDLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedCode, status)
			assert.Equal(t, tt.expectedHub, hub)
		})
	}
}

// CreateHub
type mockHubCreator struct {
	CreateHubFunc func(ctx context.Context, hub *models.Hub) error
}

func (m *mockHubCreator) CreateHub(ctx context.Context, hub *models.Hub) error {
	return m.CreateHubFunc(ctx, hub)
}

func TestCreateHubLogic(t *testing.T) {
	validTenant := uuid.New()
	hubInput := &models.Hub{Name: "Test Hub"}

	tests := []struct {
		name         string
		tenantIDStr  string
		mockFunc     func(ctx context.Context, hub *models.Hub) error
		expectedCode int
		expectErr    string
	}{
		{
			name:        "invalid tenant ID",
			tenantIDStr: "invalid-uuid",
			mockFunc: func(ctx context.Context, hub *models.Hub) error {
				return nil
			},
			expectedCode: http.StatusBadRequest,
			expectErr:    "invalid tenant_id",
		},
		{
			name:        "tenant not found",
			tenantIDStr: validTenant.String(),
			mockFunc: func(ctx context.Context, hub *models.Hub) error {
				return gorm.ErrRecordNotFound
			},
			expectedCode: http.StatusBadRequest,
			expectErr:    "tenant not found",
		},
		{
			name:        "db error",
			tenantIDStr: validTenant.String(),
			mockFunc: func(ctx context.Context, hub *models.Hub) error {
				return errors.New("db failed")
			},
			expectedCode: http.StatusInternalServerError,
			expectErr:    "failed to create hub",
		},
		{
			name:        "success",
			tenantIDStr: validTenant.String(),
			mockFunc: func(ctx context.Context, hub *models.Hub) error {
				return nil
			},
			expectedCode: http.StatusCreated,
			expectErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHubCreator{
				CreateHubFunc: tt.mockFunc,
			}

			code, err := createHubLogic(mock, tt.tenantIDStr, hubInput)

			assert.Equal(t, tt.expectedCode, code)
			if tt.expectErr != "" {
				assert.EqualError(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// DeleteHub

type mockHubDeleter struct {
	DeleteHubFunc func(ctx context.Context, id uuid.UUID) (models.Hub, error)
}

func (m *mockHubDeleter) DeleteHub(ctx context.Context, id uuid.UUID) (models.Hub, error) {
	return m.DeleteHubFunc(ctx, id)
}

func TestDeleteHubLogic(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name         string
		idStr        string
		mockFunc     func(ctx context.Context, id uuid.UUID) (models.Hub, error)
		expectedCode int
		expectErr    string
	}{
		{
			name:  "invalid UUID",
			idStr: "invalid-uuid",
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Hub, error) {
				return models.Hub{}, nil
			},
			expectedCode: http.StatusBadRequest,
			expectErr:    "invalid hub id",
		},
		{
			name:  "hub not found",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Hub, error) {
				return models.Hub{}, gorm.ErrRecordNotFound
			},
			expectedCode: http.StatusNotFound,
			expectErr:    "hub not found",
		},
		{
			name:  "success",
			idStr: validID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (models.Hub, error) {
				return models.Hub{ID: id, Name: "Test Hub"}, nil
			},
			expectedCode: http.StatusOK,
			expectErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHubDeleter{DeleteHubFunc: tt.mockFunc}
			hub, code, err := deleteHubLogic(mock, tt.idStr)

			assert.Equal(t, tt.expectedCode, code)
			if tt.expectErr != "" {
				assert.EqualError(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Test Hub", hub.Name)
			}
		})
	}
}

// UpdateHub

func TestUpdateHubLogic(t *testing.T) {
	validID := uuid.New()
	tenantID := uuid.New()
	mockHub := &models.Hub{ID: validID, Name: "Updated Hub", TenantID: tenantID}

	tests := []struct {
		name           string
		idStr          string
		input          models.Hub
		getTenantFunc  func(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
		updateFunc     func(ctx context.Context, id uuid.UUID, hub *models.Hub) error
		getFunc        func(ctx context.Context, id uuid.UUID) (*models.Hub, error)
		expectedStatus int
		expectErr      bool
	}{
		{
			name:  "invalid UUID",
			idStr: "invalid-uuid",
			input: models.Hub{},
			getTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{}, nil
			},
			updateFunc:     func(ctx context.Context, id uuid.UUID, hub *models.Hub) error { return nil },
			getFunc:        func(ctx context.Context, id uuid.UUID) (*models.Hub, error) { return mockHub, nil },
			expectedStatus: http.StatusBadRequest,
			expectErr:      true,
		},
		{
			name:  "tenant not found",
			idStr: validID.String(),
			input: models.Hub{TenantID: tenantID},
			getTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return nil, gorm.ErrRecordNotFound
			},
			updateFunc:     func(ctx context.Context, id uuid.UUID, hub *models.Hub) error { return nil },
			getFunc:        func(ctx context.Context, id uuid.UUID) (*models.Hub, error) { return mockHub, nil },
			expectedStatus: http.StatusBadRequest,
			expectErr:      true,
		},
		{
			name:  "success",
			idStr: validID.String(),
			input: models.Hub{TenantID: tenantID},
			getTenantFunc: func(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
				return &models.Tenant{}, nil
			},
			updateFunc:     func(ctx context.Context, id uuid.UUID, hub *models.Hub) error { return nil },
			getFunc:        func(ctx context.Context, id uuid.UUID) (*models.Hub, error) { return mockHub, nil },
			expectedStatus: http.StatusOK,
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockTenantService{GetTenantFunc: tt.getTenantFunc}

			result, status := updateHubLogic(
				context.Background(),
				mockService,
				tt.updateFunc,
				tt.getFunc,
				tt.idStr,
				tt.input,
			)

			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, "Updated Hub", result.Name)
			}
		})
	}
}
