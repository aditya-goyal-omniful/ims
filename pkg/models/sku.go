package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Sku struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	SkuCode   string    `gorm:"unique;not null" json:"Sku_code"`
	SellerID  uuid.UUID `gorm:"not null" json:"seller_id"`
	TenantID  uuid.UUID `gorm:"not null" json:"tenant_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func GetSkus(ctx context.Context) ([]Sku, error) {
	var skus []Sku
	if err := getDB(ctx).Find(&skus).Error; err != nil {
		return nil, err
	}
	return skus, nil
}

func GetSku(ctx context.Context, id uuid.UUID) (*Sku, error) {
	var sku Sku
	if err := getDB(ctx).First(&sku, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sku, nil
}

func CreateSku(ctx context.Context, sku *Sku) error {
	if err := getDB(ctx).Create(sku).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSku(ctx context.Context, id uuid.UUID) (Sku, error) {
	var sku Sku
	if err := getDB(ctx).First(&sku, "id = ?", id).Error; err != nil {
		return Sku{}, err
	}

	if err := getDB(ctx).Delete(&sku).Error; err != nil {
		return sku, err
	}

	return sku, nil
}

func UpdateSku(ctx context.Context, id uuid.UUID, updated *Sku) error {
	return getDB(ctx).Model(&Sku{}).Where("id = ?", id).Updates(updated).Error
}