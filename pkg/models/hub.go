package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/aditya-goyal-omniful/ims/pkg/constants"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"gorm.io/gorm"
)

type HubModel struct{}

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

// GetHubs

func (h HubModel) GetAllHubs(ctx context.Context) ([]Hub, error) {
	return GetHubs(ctx)
}

func GetHubs(ctx context.Context) ([]Hub, error) {
	var hubs []Hub
	if err := getDB(ctx).Find(&hubs).Error; err != nil {
		return nil, err
	}
	return hubs, nil
}

// GetHub

func (h HubModel) GetHub(ctx context.Context, id uuid.UUID) (*Hub, error) {
	return GetHub(ctx, id)
}

func GetHub(ctx context.Context, id uuid.UUID) (*Hub, error) {
	cacheKey := fmt.Sprintf("hub:%s", id)

	// Try to get from cache
	if cached, err := configs.RedisClient.Get(ctx, cacheKey); err == nil && cached != "" {
		var hub Hub
		if err := json.Unmarshal([]byte(cached), &hub); err == nil {
			log.Infof(i18n.Translate(ctx, "Redis cache hit for hub:"), id)
			return &hub, nil
		}
	}
	
	// Fallback to DB
	var hub Hub
	if err := getDB(ctx).First(&hub, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Store in cache
	if bytes, err := json.Marshal(hub); err == nil {
		_, err := configs.RedisClient.Set(ctx, cacheKey, string(bytes), constants.SkuCacheTTL)
		if err != nil {
			log.Infof(i18n.Translate(ctx, "Failed to set cache:"), err)
		}
	}

	return &hub, nil
}

// CreateHub

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

	// Invalidate cache
	_, _ = configs.RedisClient.Del(ctx, fmt.Sprintf("hub:%s", id))

	return hub, nil
}

func UpdateHub(ctx context.Context, id uuid.UUID, updated *Hub) error {
	err := getDB(ctx).Model(&Hub{}).Where("id = ?", id).Updates(updated).Error

	// Invalidate cache
	_, _ = configs.RedisClient.Del(ctx, fmt.Sprintf("hub:%s", id))

	return err
}