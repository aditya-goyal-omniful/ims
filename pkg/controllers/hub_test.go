package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aditya-goyal-omniful/ims/pkg/models"
	"github.com/google/uuid"
)

// GetHubs

type mockHubFetcher struct {
	GetAllHubsFunc func(ctx context.Context) ([]models.Hub, error)
}

func (m *mockHubFetcher) GetAllHubs(ctx context.Context) ([]models.Hub, error) {
	return m.GetAllHubsFunc(ctx)
}

func TestGetHubs(t *testing.T) {
	tests := []struct {
		name       string
		mockFunc   func(ctx context.Context) ([]models.Hub, error)
		wantStatus int
		wantCount  int
	}{
		{
			name: "error from DB",
			mockFunc: func(ctx context.Context) ([]models.Hub, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  0,
		},
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]models.Hub, error) {
				return []models.Hub{
					{ID: uuid.New(), Name: "Hub1"},
					{ID: uuid.New(), Name: "Hub2"},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockHubFetcher{
				GetAllHubsFunc: tt.mockFunc,
			}

			result, status := getHubsLogic(service)

			if status != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, status)
			}
			if len(result) != tt.wantCount {
				t.Errorf("expected %d hubs, got %d", tt.wantCount, len(result))
			}
		})
	}
}

// GetHub

type mockHubService struct {
	GetHubFunc func(ctx context.Context, id uuid.UUID) (*models.Hub, error)
}

func (m *mockHubService) GetHub(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
	return m.GetHubFunc(ctx, id)
}

func TestGetHubByID(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name         string
		inputID      string
		mockFunc     func(ctx context.Context, id uuid.UUID) (*models.Hub, error)
		wantStatus   int
		wantNilHub   bool
		wantIDMatch  bool
	}{
		{
			name:       "invalid UUID",
			inputID:    "not-a-uuid",
			mockFunc:   nil, // should not be called
			wantStatus: http.StatusBadRequest,
			wantNilHub: true,
		},
		{
			name:    "DB error",
			inputID: testID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
				return nil, errors.New("DB error")
			},
			wantStatus: http.StatusInternalServerError,
			wantNilHub: true,
		},
		{
			name:    "success",
			inputID: testID.String(),
			mockFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
				return &models.Hub{ID: id}, nil
			},
			wantStatus:  http.StatusOK,
			wantNilHub:  false,
			wantIDMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mockHubService{
				GetHubFunc: func(ctx context.Context, id uuid.UUID) (*models.Hub, error) {
					if tt.mockFunc != nil {
						return tt.mockFunc(ctx, id)
					}
					t.Fatal("mockFunc should not be called")
					return nil, nil
				},
			}

			result, status := getHubByIDLogic(service, tt.inputID)

			if status != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, status)
			}
			if tt.wantNilHub && result != nil {
				t.Errorf("expected nil hub, got %+v", result)
			}
			if !tt.wantNilHub && result == nil {
				t.Errorf("expected hub, got nil")
			}
			if tt.wantIDMatch && result != nil && result.ID.String() != tt.inputID {
				t.Errorf("expected hub ID %s, got %s", tt.inputID, result.ID.String())
			}
		})
	}
}

// CreateHub

