package models

import (
	"context"
	"time"

	"github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hub struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Location  string    `json:"location"`
	TenantID  uuid.UUID `gorm:"not null" json:"tenant_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func getDB(ctx context.Context) *gorm.DB {
	return configs.GetDB().GetMasterDB(ctx)
}

func GetHubs(ctx context.Context) ([]Hub, error) {
	var hubs []Hub
	if err := getDB(ctx).Find(&hubs).Error; err != nil {
		return nil, err
	}
	return hubs, nil
}

func GetHub(ctx context.Context, id uuid.UUID) (*Hub, error) {
	var hub Hub
	if err := getDB(ctx).First(&hub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &hub, nil
}

func CreateHub(ctx context.Context, hub *Hub) error {
	// Check if tenant exists before creating hub
	_, err := GetTenant(ctx, hub.TenantID)
	if err != nil {
		return err // This will be a gorm.ErrRecordNotFound if tenant doesn't exist
	}

	if err := getDB(ctx).Create(hub).Error; err != nil {
		return err
	}
	return nil
}

func DeleteHub(ctx context.Context, id uuid.UUID) (Hub, error) {
	var hub Hub
	if err := getDB(ctx).First(&hub, "id = ?", id).Error; err != nil {
		return Hub{}, err
	}

	if err := getDB(ctx).Delete(&hub).Error; err != nil {
		return Hub{}, err
	}

	return hub, nil
}

func UpdateHub(ctx context.Context, id uuid.UUID, updated *Hub) error {
	return getDB(ctx).Model(&Hub{}).Where("id = ?", id).Updates(updated).Error
}