package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	HubID     uuid.UUID `gorm:"type:uuid;not null" json:"hub_id"`
	SkuID     uuid.UUID `gorm:"type:uuid;not null" json:"sku_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func GetInventories(ctx context.Context) ([]Inventory, error) {
	var inventories []Inventory
	if err := getDB(ctx).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

func GetInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	var inventory Inventory
	if err := getDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

func CreateInventory(ctx context.Context, inventory *Inventory) error {
	// Check if tenant exists before creating inventory
	_, err := GetTenant(ctx, inventory.TenantID)
	if err != nil {
		return err // This will be a gorm.ErrRecordNotFound if tenant doesn't exist
	}

	if err := getDB(ctx).Create(inventory).Error; err != nil {
		return err
	}
	return nil
}

func DeleteInventory(ctx context.Context, id uuid.UUID) (Inventory, error) {
	var inventory Inventory
	if err := getDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return Inventory{}, err
	}

	if err := getDB(ctx).Delete(&inventory).Error; err != nil {
		return inventory, err
	}

	return inventory, nil
}

func UpdateInventory(ctx context.Context, id uuid.UUID, updated *Inventory) error {
	return getDB(ctx).Model(&Inventory{}).Where("id = ?", id).Updates(updated).Error
}