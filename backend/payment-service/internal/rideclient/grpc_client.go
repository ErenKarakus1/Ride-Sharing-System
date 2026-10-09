package rideclient

import (
	"context"

	ridev1 "github.com/ErenKarakus1/Ride-Sharing-System/backend/proto/gen/go/ride/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const internalTokenHeader = "x-internal-service-token"

type Ride struct {
	ID      string
	RiderID string
	Status  string
}

type Client interface {
	GetRide(ctx context.Context, id string) (Ride, error)
}

type GRPCClient struct {
	client        ridev1.RideServiceClient
	internalToken string
}

func NewGRPCClient(conn *grpc.ClientConn, internalToken string) *GRPCClient {
	return &GRPCClient{
		client:        ridev1.NewRideServiceClient(conn),
		internalToken: internalToken,
	}
}

func (c *GRPCClient) GetRide(ctx context.Context, id string) (Ride, error) {
	if c.internalToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, internalTokenHeader, c.internalToken)
	}

	response, err := c.client.GetRide(ctx, &ridev1.GetRideRequest{Id: id})
	if err != nil {
		return Ride{}, err
	}

	return Ride{
		ID:      response.GetId(),
		RiderID: response.GetRiderId(),
		Status:  response.GetStatus(),
	}, nil
}
