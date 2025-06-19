package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null;unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetTenants(ctx context.Context) ([]Tenant, error) {
	var tenants []Tenant
	if err := getDB(ctx).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func GetTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	var tenant Tenant
	if err := getDB(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func CreateTenant(ctx context.Context, tenant *Tenant) error {
	if err := getDB(ctx).Create(tenant).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTenant(ctx context.Context, id uuid.UUID) (Tenant, error) {
	var tenant Tenant
	if err := getDB(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return Tenant{}, err
	}

	if err := getDB(ctx).Delete(&tenant).Error; err != nil {
		return tenant, err
	}

	return tenant, nil
}

func UpdateTenant(ctx context.Context, id uuid.UUID, updated *Tenant) error {
	return getDB(ctx).Model(&Tenant{}).Where("id = ?", id).Updates(updated).Error
}