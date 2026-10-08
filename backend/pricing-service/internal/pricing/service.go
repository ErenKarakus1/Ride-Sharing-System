package pricing

import "math"

type Service struct {
	baseFare  float64
	perKM     float64
	perMinute float64
}

func NewService(baseFare float64, perKM float64, perMinute float64) *Service {
	return &Service{
		baseFare:  baseFare,
		perKM:     perKM,
		perMinute: perMinute,
	}
}

func (s *Service) Estimate(request EstimateRequest) EstimateResponse {
	distanceKM := haversineKM(
		request.Pickup.Latitude,
		request.Pickup.Longitude,
		request.Dropoff.Latitude,
		request.Dropoff.Longitude,
	)
	durationMinutes := distanceKM * 2.4
	total := s.baseFare + distanceKM*s.perKM + durationMinutes*s.perMinute

	return EstimateResponse{
		DistanceKM:      round(distanceKM, 2),
		DurationMinutes: round(durationMinutes, 0),
		Currency:        "TRY",
		Amount:          round(total, 2),
	}
}

func haversineKM(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const earthRadiusKM = 6371
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	return earthRadiusKM * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func degreesToRadians(value float64) float64 {
	return value * math.Pi / 180
}

func round(value float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Round(value*pow) / pow
}
