package location

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const driverLocationsKey = "driver_locations"
const availableDriversKey = "available_drivers"
const driverAvailabilityTTL = 15 * time.Minute

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
	driverID = strings.TrimSpace(driverID)
	pipe := r.client.TxPipeline()
	pipe.SAdd(ctx, availableDriversKey, driverID)
	pipe.Set(ctx, driverAvailabilityKey(driverID), "1", driverAvailabilityTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisRepository) SetDriverUnavailable(ctx context.Context, driverID string) error {
	driverID = strings.TrimSpace(driverID)
	pipe := r.client.TxPipeline()
	pipe.SRem(ctx, availableDriversKey, driverID)
	pipe.Del(ctx, driverAvailabilityKey(driverID))
	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisRepository) ClaimDriver(ctx context.Context, driverID string) error {
	driverID = strings.TrimSpace(driverID)
	if driverID == "" {
		return ErrMissingDriver
	}

	claimed, err := r.client.Eval(ctx, `
		if redis.call("SISMEMBER", KEYS[1], ARGV[1]) == 0 then
			return 0
		end
		if redis.call("EXISTS", KEYS[2]) == 0 then
			redis.call("SREM", KEYS[1], ARGV[1])
			return 0
		end
		redis.call("SREM", KEYS[1], ARGV[1])
		redis.call("DEL", KEYS[2])
		return 1
	`, []string{availableDriversKey, driverAvailabilityKey(driverID)}, driverID).Int()
	if err != nil {
		return err
	}
	if claimed == 0 {
		return ErrDriverUnavailable
	}

	return nil
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
		available, err := r.driverAvailable(ctx, result.Name)
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

func driverAvailabilityKey(driverID string) string {
	return fmt.Sprintf("driver:%s:available", driverID)
}

func (r *RedisRepository) driverAvailable(ctx context.Context, driverID string) (bool, error) {
	available, err := r.client.SIsMember(ctx, availableDriversKey, driverID).Result()
	if err != nil || !available {
		return available, err
	}

	live, err := r.client.Exists(ctx, driverAvailabilityKey(driverID)).Result()
	if err != nil {
		return false, err
	}
	if live == 0 {
		_ = r.client.SRem(ctx, availableDriversKey, driverID).Err()
		return false, nil
	}

	return true, nil
}
