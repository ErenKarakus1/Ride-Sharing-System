package matching

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Match(ctx *gin.Context) {
	var request MatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := h.service.MatchAndAssign(ctx.Request.Context(), request)
	if errors.Is(err, ErrNoDriversAvailable) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no drivers available"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to match ride"})
		return
	}

	ctx.JSON(http.StatusOK, match)
}

type MatchRequest struct {
	RideID   string   `json:"ride_id" binding:"required"`
	Pickup   Location `json:"pickup" binding:"required"`
	RadiusKM float64  `json:"radius_km"`
	Limit    int      `json:"limit"`
}

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type MatchResponse struct {
	RideID    string  `json:"ride_id"`
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
