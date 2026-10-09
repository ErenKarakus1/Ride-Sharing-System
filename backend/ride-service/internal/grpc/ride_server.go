package grpc

import (
	"context"
	"errors"
	"time"

	ridev1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/ride/v1"
	"github.com/ErenKarakus1/Ride-Sharing-System/backend/ride-service/internal/ride"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RideServer struct {
	ridev1.UnimplementedRideServiceServer
	service *ride.Service
}

func NewRideServer(service *ride.Service) *RideServer {
	return &RideServer{
		service: service,
	}
}

func (s *RideServer) GetRide(ctx context.Context, request *ridev1.GetRideRequest) (*ridev1.Ride, error) {
	found, err := s.service.Get(ctx, request.GetId())
	if errors.Is(err, ride.ErrRideNotFound) {
		return nil, status.Error(codes.NotFound, "ride not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get ride")
	}

	return rideToProto(found), nil
}

func (s *RideServer) ListRiderRides(ctx context.Context, request *ridev1.ListRiderRidesRequest) (*ridev1.ListRiderRidesResponse, error) {
	rides, err := s.service.ListByRider(ctx, request.GetRiderId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list rides")
	}

	response := &ridev1.ListRiderRidesResponse{
		Rides: make([]*ridev1.Ride, 0, len(rides)),
	}
	for _, item := range rides {
		response.Rides = append(response.Rides, rideToProto(item))
	}

	return response, nil
}

func rideToProto(value ride.Ride) *ridev1.Ride {
	driverID := ""
	if value.DriverID != nil {
		driverID = *value.DriverID
	}

	return &ridev1.Ride{
		Id:       value.ID,
		RiderId:  value.RiderID,
		DriverId: driverID,
		Pickup: &ridev1.Location{
			Latitude:  value.Pickup.Latitude,
			Longitude: value.Pickup.Longitude,
			Address:   value.Pickup.Address,
		},
		Dropoff: &ridev1.Location{
			Latitude:  value.Dropoff.Latitude,
			Longitude: value.Dropoff.Longitude,
			Address:   value.Dropoff.Address,
		},
		Status:    string(value.Status),
		CreatedAt: formatTime(value.CreatedAt),
		UpdatedAt: formatTime(value.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
