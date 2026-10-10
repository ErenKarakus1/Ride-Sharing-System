package location

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UpdateDriverLocation(ctx *gin.Context) {
	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	var request UpdateDriverLocationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateDriverLocation(ctx.Request.Context(), driverID(ctx), request)
	h.writeEmptyResult(ctx, err)
}

func (h *Handler) SetDriverAvailable(ctx *gin.Context) {
	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	h.writeEmptyResult(ctx, h.service.SetDriverAvailable(ctx.Request.Context(), driverID(ctx)))
}

func (h *Handler) SetDriverUnavailable(ctx *gin.Context) {
	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	h.writeEmptyResult(ctx, h.service.SetDriverUnavailable(ctx.Request.Context(), driverID(ctx)))
}

func (h *Handler) ClaimDriver(ctx *gin.Context) {
	h.writeEmptyResult(ctx, h.service.ClaimDriver(ctx.Request.Context(), driverID(ctx)))
}

func (h *Handler) NearbyDrivers(ctx *gin.Context) {
	request := NearbyDriversRequest{
		Latitude:  floatQuery(ctx, "latitude"),
		Longitude: floatQuery(ctx, "longitude"),
		RadiusKM:  floatQuery(ctx, "radius_km"),
		Limit:     intQuery(ctx, "limit"),
	}

	drivers, err := h.service.NearbyDrivers(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find nearby drivers"})
		return
	}

	ctx.JSON(http.StatusOK, drivers)
}

func (h *Handler) writeEmptyResult(ctx *gin.Context, err error) {
	if errors.Is(err, ErrMissingDriver) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing driver"})
		return
	}
	if errors.Is(err, ErrDriverUnavailable) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "driver unavailable"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update driver location state"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func floatQuery(ctx *gin.Context, key string) float64 {
	value, _ := strconv.ParseFloat(ctx.Query(key), 64)
	return value
}

func intQuery(ctx *gin.Context, key string) int {
	value, _ := strconv.Atoi(ctx.Query(key))
	return value
}

type UpdateDriverLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type NearbyDriversRequest struct {
	Latitude  float64
	Longitude float64
	RadiusKM  float64
	Limit     int
}

func driverID(ctx *gin.Context) string {
	if id := ctx.GetHeader("X-User-ID"); id != "" {
		return id
	}

	return ctx.Param("id")
}
