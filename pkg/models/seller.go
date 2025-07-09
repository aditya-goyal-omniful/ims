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

type SellerModel struct{}

// GetSellers

func (s SellerModel) GetSellers(ctx context.Context) ([]Seller, error) {
	return GetSellers(ctx)
}

func GetSellers(ctx context.Context) ([]Seller, error) {
	var sellers []Seller
	if err := getDB(ctx).Find(&sellers).Error; err != nil {
		return nil, err
	}
	return sellers, nil
}

// GetSellerByID

func (s SellerModel) GetSeller(ctx context.Context, id uuid.UUID) (*Seller, error) {
	return GetSeller(ctx, id)
}

func GetSeller(ctx context.Context, id uuid.UUID) (*Seller, error) {
	var seller Seller
	if err := getDB(ctx).First(&seller, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}

// CreateSeller

func (s SellerModel) CreateSeller(ctx context.Context, seller *Seller) error {
	return CreateSeller(ctx, seller)
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

// DeleteSeller

func (s SellerModel) DeleteSeller(ctx context.Context, id uuid.UUID) (*Seller, error) {
	seller, err := GetSeller(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := getDB(ctx).Delete(seller).Error; err != nil {
		return nil, err
	}

	return seller, nil
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

// UpdateSeller

func (s SellerModel) UpdateSeller(ctx context.Context, id uuid.UUID, seller *Seller) error {
	return UpdateSeller(ctx, id, seller)
}

func (s SellerModel) GetTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return GetTenant(ctx, id)
}

func UpdateSeller(ctx context.Context, id uuid.UUID, updated *Seller) error {
	return getDB(ctx).Model(&Seller{}).Where("id = ?", id).Updates(updated).Error
}