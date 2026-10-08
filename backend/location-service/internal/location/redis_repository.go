package location

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const driverLocationsKey = "driver_locations"
const availableDriversKey = "available_drivers"

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) UpdateDriverLocation(ctx context.Context, location DriverLocation) error {
	driverID := strings.TrimSpace(location.DriverID)
	pipe := r.client.TxPipeline()
	pipe.GeoAdd(ctx, driverLocationsKey, &redis.GeoLocation{
		Name:      driverID,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	})
	pipe.HSet(ctx, driverLocationHashKey(driverID), map[string]any{
		"latitude":  strconv.FormatFloat(location.Latitude, 'f', -1, 64),
		"longitude": strconv.FormatFloat(location.Longitude, 'f', -1, 64),
	})
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisRepository) SetDriverAvailable(ctx context.Context, driverID string) error {
	return r.client.SAdd(ctx, availableDriversKey, strings.TrimSpace(driverID)).Err()
}

func (r *RedisRepository) SetDriverUnavailable(ctx context.Context, driverID string) error {
	return r.client.SRem(ctx, availableDriversKey, strings.TrimSpace(driverID)).Err()
}

func (r *RedisRepository) NearbyDrivers(ctx context.Context, latitude float64, longitude float64, radiusKM float64, limit int) ([]DriverLocation, error) {
	results, err := r.client.GeoSearchLocation(ctx, driverLocationsKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  longitude,
			Latitude:   latitude,
			Radius:     radiusKM,
			RadiusUnit: "km",
			Sort:       "ASC",
			Count:      limit,
		},
		WithCoord: true,
	}).Result()
	if err != nil {
		return nil, err
	}

	drivers := make([]DriverLocation, 0, len(results))
	for _, result := range results {
		available, err := r.client.SIsMember(ctx, availableDriversKey, result.Name).Result()
		if err != nil {
			return nil, err
		}
		if !available {
			continue
		}

		drivers = append(drivers, DriverLocation{
			DriverID:  result.Name,
			Latitude:  result.Latitude,
			Longitude: result.Longitude,
		})
	}

	return drivers, nil
}

func driverLocationHashKey(driverID string) string {
	return fmt.Sprintf("driver:%s:location", driverID)
}
