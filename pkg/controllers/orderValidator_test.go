package controllers

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/omniful/go_commons/http"
	"github.com/stretchr/testify/assert"
)

type mockValidator struct {
	ValidateHubAndSkuFunc func(ctx context.Context, hubID, skuID uuid.UUID) (bool, error)
}

func (m *mockValidator) ValidateHubAndSku(ctx context.Context, hubID, skuID uuid.UUID) (bool, error) {
	return m.ValidateHubAndSkuFunc(ctx, hubID, skuID)
}

func TestValidateOrderLogic(t *testing.T) {
	validHubID := uuid.New()
	validSkuID := uuid.New()

	tests := []struct {
		name           string
		hubIDStr       string
		skuIDStr       string
		mockFunc       func(ctx context.Context, hubID, skuID uuid.UUID) (bool, error)
		expectedValid  bool
		expectedStatus int
		expectErr      bool
	}{
		{
			name:           "invalid hub ID",
			hubIDStr:       "bad-id",
			skuIDStr:       validSkuID.String(),
			expectedValid:  false,
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:           "invalid sku ID",
			hubIDStr:       validHubID.String(),
			skuIDStr:       "bad-id",
			expectedValid:  false,
			expectedStatus: int(http.StatusBadRequest),
			expectErr:      true,
		},
		{
			name:     "validation error",
			hubIDStr: validHubID.String(),
			skuIDStr: validSkuID.String(),
			mockFunc: func(ctx context.Context, hubID, skuID uuid.UUID) (bool, error) {
				return false, errors.New("db error")
			},
			expectedValid:  false,
			expectedStatus: int(http.StatusInternalServerError),
			expectErr:      true,
		},
		{
			name:     "successful validation",
			hubIDStr: validHubID.String(),
			skuIDStr: validSkuID.String(),
			mockFunc: func(ctx context.Context, hubID, skuID uuid.UUID) (bool, error) {
				return true, nil
			},
			expectedValid:  true,
			expectedStatus: int(http.StatusOK),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockValidator{
				ValidateHubAndSkuFunc: tt.mockFunc,
			}

			isValid, status, err := validateOrderLogic(mock, tt.hubIDStr, tt.skuIDStr)
			assert.Equal(t, tt.expectedValid, isValid)
			assert.Equal(t, tt.expectedStatus, status)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
