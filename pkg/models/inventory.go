package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Location  string    `json:"location"`
	TenantID  uuid.UUID `gorm:"not null" json:"tenant_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func GetInventories(ctx context.Context) ([]Inventory, error) {
	var inventories []Inventory
	if err := db.GetMasterDB(ctx).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

func GetInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	var inventory Inventory
	if err := db.GetMasterDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

func CreateInventory(ctx context.Context, inventory *Inventory) error {
	if err := db.GetMasterDB(ctx).Create(inventory).Error; err != nil {
		return err
	}
	return nil
}

func DeleteInventory(ctx context.Context, id uuid.UUID) (Inventory, error) {
	var inventory Inventory
	if err := db.GetMasterDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return Inventory{}, err
	}

	if err := db.GetMasterDB(ctx).Delete(&inventory).Error; err != nil {
		return inventory, err
	}

	return inventory, nil
}

func UpdateInventory(ctx context.Context, id uuid.UUID, updated *Inventory) error {
	return db.GetMasterDB(ctx).Model(&Inventory{}).Where("id = ?", id).Updates(updated).Error
}