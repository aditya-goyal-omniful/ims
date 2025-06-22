package models

import (
	"context"
	"fmt"

	"github.com/aditya-goyal-omniful/ims/pkg/configs"
	"github.com/aditya-goyal-omniful/ims/pkg/constants"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/log"
)

func ValidateHubAndSku(ctx context.Context, hubID, skuID uuid.UUID) (bool, error) {
	hubKey := fmt.Sprintf("hub_valid:%s", hubID)
	skuKey := fmt.Sprintf("sku_valid:%s", skuID)

	var hubValid bool
	if cached, err := configs.RedisClient.Get(ctx, hubKey); err == nil && cached == "true" {
		log.Infof("Hub %s found in Redis.", hubID)
		hubValid = true
	} else {
		hub, err := GetHub(ctx, hubID)
		if err != nil || hub == nil {
			log.Warnf("Hub %s not found in DB.", hubID)
			return false, err
		}
		hubValid = true
		_, _ = configs.RedisClient.Set(ctx, hubKey, "true", constants.RedisCacheTTL)
		log.Infof("Hub %s cached in Redis.", hubID)
	}

	var skuValid bool
	if cached, err := configs.RedisClient.Get(ctx, skuKey); err == nil && cached == "true" {
		log.Infof("SKU %s found in Redis.", skuID)
		skuValid = true
	} else {
		sku, err := GetSku(ctx, skuID)
		if err != nil || sku == nil {
			log.Warnf("SKU %s not found in DB.", skuID)
			return false, err
		}
		skuValid = true
		_, _ = configs.RedisClient.Set(ctx, skuKey, "true", constants.RedisCacheTTL)
		log.Infof("SKU %s cached in Redis.", skuID)
	}

	return hubValid && skuValid, nil
}