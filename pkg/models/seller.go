package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	TenantID  uuid.UUID `gorm:"not null" json:"tenant_id"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func GetSellers(ctx context.Context) ([]Seller, error) {
	var sellers []Seller
	if err := getDB(ctx).Find(&sellers).Error; err != nil {
		return nil, err
	}
	return sellers, nil
}

func GetSeller(ctx context.Context, id uuid.UUID) (*Seller, error) {
	var seller Seller
	if err := getDB(ctx).First(&seller, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}

func CreateSeller(ctx context.Context, seller *Seller) error {
	// Check if tenant exists before creating seller
	_, err := GetTenant(ctx, seller.TenantID)
	if err != nil {
		return err // This will be a gorm.ErrRecordNotFound if tenant doesn't exist
	}

	if err := getDB(ctx).Create(seller).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSeller(ctx context.Context, id uuid.UUID) (Seller, error) {
	var seller Seller
	if err := getDB(ctx).First(&seller, "id = ?", id).Error; err != nil {
		return Seller{}, err
	}

	if err := getDB(ctx).Delete(&seller).Error; err != nil {
		return seller, err
	}

	return seller, nil
}

func UpdateSeller(ctx context.Context, id uuid.UUID, updated *Seller) error {
	return getDB(ctx).Model(&Seller{}).Where("id = ?", id).Updates(updated).Error
}