package ride

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx *gin.Context) {
	if ctx.GetHeader("X-User-Role") != "rider" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "rider role required"})
		return
	}

	var request CreateRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request.RiderID = ctx.GetHeader("X-User-ID")

	created, err := h.service.Create(ctx.Request.Context(), request)
	if errors.Is(err, ErrMissingRider) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing rider"})
		return
	}
	if errors.Is(err, ErrInvalidLocation) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid location"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ride"})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (h *Handler) Get(ctx *gin.Context) {
	if !validUUIDParam(ctx, "id") {
		return
	}

	found, err := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if errors.Is(err, ErrRideNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ride not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ride"})
		return
	}

	ctx.JSON(http.StatusOK, found)
}

func (h *Handler) ListByRider(ctx *gin.Context) {
	rides, err := h.service.ListByRider(ctx.Request.Context(), ctx.Query("rider_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rides"})
		return
	}

	ctx.JSON(http.StatusOK, rides)
}

func (h *Handler) Accept(ctx *gin.Context) {
	if !validUUIDParam(ctx, "id") {
		return
	}

	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	var request AcceptRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.DriverID == "" {
		request.DriverID = ctx.GetHeader("X-User-ID")
	}

	ride, err := h.service.Accept(ctx.Request.Context(), ctx.Param("id"), request.DriverID)
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Start(ctx *gin.Context) {
	if !validUUIDParam(ctx, "id") {
		return
	}

	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	ride, err := h.service.Start(ctx.Request.Context(), ctx.Param("id"), ctx.GetHeader("X-User-ID"))
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Complete(ctx *gin.Context) {
	if !validUUIDParam(ctx, "id") {
		return
	}

	if ctx.GetHeader("X-User-Role") != "driver" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "driver role required"})
		return
	}

	ride, err := h.service.Complete(ctx.Request.Context(), ctx.Param("id"), ctx.GetHeader("X-User-ID"))
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Cancel(ctx *gin.Context) {
	if !validUUIDParam(ctx, "id") {
		return
	}

	ride, err := h.service.Cancel(ctx.Request.Context(), ctx.Param("id"), ctx.GetHeader("X-User-ID"))
	h.writeRideResult(ctx, ride, err)
}

func validUUIDParam(ctx *gin.Context, name string) bool {
	if _, err := uuid.Parse(strings.TrimSpace(ctx.Param(name))); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return false
	}

	return true
}

func (h *Handler) writeRideResult(ctx *gin.Context, ride Ride, err error) {
	if errors.Is(err, ErrRideNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ride not found"})
		return
	}
	if errors.Is(err, ErrInvalidTransition) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "invalid ride status transition"})
		return
	}
	if errors.Is(err, ErrUnauthorizedRideAction) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized ride action"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ride"})
		return
	}

	ctx.JSON(http.StatusOK, ride)
}

type CreateRideRequest struct {
	RiderID string   `json:"rider_id"`
	Pickup  Location `json:"pickup" binding:"required"`
	Dropoff Location `json:"dropoff" binding:"required"`
}

type AcceptRideRequest struct {
	DriverID string `json:"driver_id"`
}
