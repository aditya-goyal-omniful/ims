package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
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

type InventoryView struct {
	SkuID    uuid.UUID `json:"sku_id"`
	SkuCode  string    `json:"sku_code"`
	SkuName  string    `json:"sku_name"`
	Quantity int       `json:"quantity"`
}

type InventoryModel struct{}

// GetInventories

func (i InventoryModel) GetInventories(ctx context.Context) ([]Inventory, error) {
	return GetInventories(ctx)
}

func GetInventories(ctx context.Context) ([]Inventory, error) {
	var inventories []Inventory
	if err := getDB(ctx).Find(&inventories).Error; err != nil {
		return nil, err
	}
	return inventories, nil
}

// GetInventoryByID

func (i InventoryModel) GetInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	return GetInventory(ctx, id)
}

func GetInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	var inventory Inventory
	if err := getDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

// CreateInventory

func (i InventoryModel) CreateInventory(ctx context.Context, inv *Inventory) error {
	return CreateInventory(ctx, inv)
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

// DeleteInventory

func (i InventoryModel) DeleteInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	return DeleteInventory(ctx, id)
}

func DeleteInventory(ctx context.Context, id uuid.UUID) (*Inventory, error) {
	var inventory Inventory
	if err := getDB(ctx).First(&inventory, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := getDB(ctx).Delete(&inventory).Error; err != nil {
		return nil, err
	}

	return &inventory, nil
}

// UpdateInventory

func (i InventoryModel) UpdateInventory(ctx context.Context, id uuid.UUID, updated *Inventory) error {
	return UpdateInventory(ctx, id, updated)
}

func UpdateInventory(ctx context.Context, id uuid.UUID, updated *Inventory) error {
	return getDB(ctx).Model(&Inventory{}).Where("id = ?", id).Updates(updated).Error
}

// UpsertInventory

func (i InventoryModel) UpsertInventory(ctx context.Context, inv *Inventory) error {
	return UpsertInventory(ctx, inv)
}

func UpsertInventory(ctx context.Context, inventory *Inventory) error {
	// Validate tenant exists
	if _, err := GetTenant(ctx, inventory.TenantID); err != nil {
		return err
	}

	db := getDB(ctx)

	// Atomic UPSERT: (sku_id, hub_id) must be unique for this to work properly
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sku_id"}, {Name: "hub_id"}}, // conflict target
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(inventory).Error
}

// GetInventoryWithDefaults

func (i InventoryModel) GetInventoryWithDefaults(ctx context.Context, tenantID, hubID uuid.UUID) ([]InventoryView, error) {
	return GetInventoryWithDefaults(ctx, tenantID, hubID)
}

func GetInventoryWithDefaults(ctx context.Context, tenantID, hubID uuid.UUID) ([]InventoryView, error) {
	var result []InventoryView
	db := getDB(ctx)

	err := db.Raw(`
		SELECT 
			s.id AS sku_id,
			s.sku_code,
			s.name AS sku_name,
			COALESCE(i.quantity, 0) AS quantity
		FROM skus s
		LEFT JOIN inventories i 
			ON s.id = i.sku_id AND i.hub_id = ? AND i.tenant_id = s.tenant_id
		WHERE s.tenant_id = ?
	`, hubID, tenantID).Scan(&result).Error

	return result, err
}

// GetInventoryBySkuHub

func GetInventoryBySkuHub(ctx context.Context, skuID, hubID uuid.UUID) (*Inventory, error) {
	var inv Inventory
	err := getDB(ctx).Where("sku_id = ? AND hub_id = ?", skuID, hubID).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// UpdateInventoryQuantity

func (m InventoryModel) GetInventoryBySkuHub(ctx context.Context, skuID, hubID uuid.UUID) (*Inventory, error) {
	return GetInventoryBySkuHub(ctx, skuID, hubID)
}

func (m InventoryModel) UpdateInventoryQuantity(ctx context.Context, invID uuid.UUID, newQty int) error {
	return UpdateInventoryQuantity(ctx, invID, newQty)
}

func UpdateInventoryQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	return getDB(ctx).Model(&Inventory{}).Where("id = ?", id).Update("quantity", quantity).Error
}
