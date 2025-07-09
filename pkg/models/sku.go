package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/aditya-goyal-omniful/ims/pkg/constants"
	"github.com/google/uuid"
)

type Sku struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	SkuCode   string    `gorm:"unique;not null" json:"sku_code"`
	SellerID  uuid.UUID `gorm:"not null" json:"seller_id"`
	TenantID  uuid.UUID `gorm:"not null" json:"tenant_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type SKUModel struct{}

// GetSkus

func GetSkus(ctx context.Context) ([]Sku, error) {
	var skus []Sku
	if err := getDB(ctx).Find(&skus).Error; err != nil {
		return nil, err
	}
	return skus, nil
}

// GetSkuByID

func (s SKUModel) GetSku(ctx context.Context, id uuid.UUID) (*Sku, error) {
	return GetSku(ctx, id)
}

func GetSku(ctx context.Context, id uuid.UUID) (*Sku, error) {
	cacheKey := fmt.Sprintf("sku:%s", id)

	// Try to get from cache
	if cached, err := configs.RedisClient.Get(ctx, cacheKey); err == nil && cached != "" {
		var sku Sku
		if err := json.Unmarshal([]byte(cached), &sku); err == nil {
			fmt.Println("Redis cache hit for sku:", id)
			return &sku, nil
		}
	}
	
	// Fallback to DB
	var sku Sku
	if err := getDB(ctx).First(&sku, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Store in cache
	if bytes, err := json.Marshal(sku); err == nil {
		_, _ = configs.RedisClient.Set(ctx, cacheKey, string(bytes), constants.SkuCacheTTL)
	}

	return &sku, nil
}

// CreateSku

func (s SKUModel) CreateSku(ctx context.Context, sku *Sku) error {
	return CreateSku(ctx, sku)
}

func CreateSku(ctx context.Context, sku *Sku) error {
	// Check if tenant exists before creating sku
	_, err := GetTenant(ctx, sku.TenantID)
	if err != nil {
		return err // This will be a gorm.ErrRecordNotFound if tenant doesn't exist
	}

	// Check if seller exists before creating sku
	_, err = GetSeller(ctx, sku.SellerID)
	if err != nil {
		return err // This will be a gorm.ErrRecordNotFound if seller doesn't exist
	}

	if err := getDB(ctx).Create(sku).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSku

func (s SKUModel) DeleteSku(ctx context.Context, id uuid.UUID) (*Sku, error) {
	return DeleteSku(ctx, id)
}

func DeleteSku(ctx context.Context, id uuid.UUID) (*Sku, error) {
	var sku Sku
	if err := getDB(ctx).First(&sku, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := getDB(ctx).Delete(&sku).Error; err != nil {
		return nil, err
	}

	// Invalidate cache
	_, _ = configs.RedisClient.Del(ctx, fmt.Sprintf("sku:%s", id))

	return &sku, nil
}

// UpdateSku

func (s SKUModel) UpdateSku(ctx context.Context, id uuid.UUID, updated *Sku) error {
	return UpdateSku(ctx, id, updated)
}

func (s SKUModel) GetTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return GetTenant(ctx, id)
}

func UpdateSku(ctx context.Context, id uuid.UUID, updated *Sku) error {
	err := getDB(ctx).Model(&Sku{}).Where("id = ?", id).Updates(updated).Error

	// Invalidate cache
	_, _ = configs.RedisClient.Del(ctx, fmt.Sprintf("sku:%s", id))

	return err
}

// GetFilteredSkus

func (s SKUModel) GetFilteredSkus(ctx context.Context, tenantID, sellerID uuid.UUID, skuCodes []string) ([]Sku, error) {
	return GetFilteredSkus(ctx, tenantID, sellerID, skuCodes)
}

func GetFilteredSkus(ctx context.Context, tenantID uuid.UUID, sellerID uuid.UUID, skuCodes []string) ([]Sku, error) {
	db := getDB(ctx)
	query := db.Model(&Sku{}).Where("tenant_id = ?", tenantID)

	if sellerID != uuid.Nil {
		query = query.Where("seller_id = ?", sellerID)
	}

	if len(skuCodes) > 0 {
		query = query.Where("sku_code IN ?", skuCodes)
	}

	var skus []Sku
	if err := query.Find(&skus).Error; err != nil {
		return nil, err
	}

	return skus, nil
}