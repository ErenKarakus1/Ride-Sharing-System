package pricing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Estimate(ctx *gin.Context) {
	var request EstimateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, h.service.Estimate(request))
}

type EstimateRequest struct {
	Pickup  Location `json:"pickup" binding:"required"`
	Dropoff Location `json:"dropoff" binding:"required"`
}

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type EstimateResponse struct {
	DistanceKM      float64 `json:"distance_km"`
	DurationMinutes float64 `json:"duration_minutes"`
	Currency        string  `json:"currency"`
	Amount          float64 `json:"amount"`
}
