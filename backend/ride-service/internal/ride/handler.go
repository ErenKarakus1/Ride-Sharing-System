package ride

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

func (h *Handler) Create(ctx *gin.Context) {
	var request CreateRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.Create(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ride"})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (h *Handler) Get(ctx *gin.Context) {
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
	var request AcceptRideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ride, err := h.service.Accept(ctx.Request.Context(), ctx.Param("id"), request.DriverID)
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Start(ctx *gin.Context) {
	ride, err := h.service.Start(ctx.Request.Context(), ctx.Param("id"))
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Complete(ctx *gin.Context) {
	ride, err := h.service.Complete(ctx.Request.Context(), ctx.Param("id"))
	h.writeRideResult(ctx, ride, err)
}

func (h *Handler) Cancel(ctx *gin.Context) {
	ride, err := h.service.Cancel(ctx.Request.Context(), ctx.Param("id"))
	h.writeRideResult(ctx, ride, err)
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ride"})
		return
	}

	ctx.JSON(http.StatusOK, ride)
}

type CreateRideRequest struct {
	RiderID string   `json:"rider_id" binding:"required"`
	Pickup  Location `json:"pickup" binding:"required"`
	Dropoff Location `json:"dropoff" binding:"required"`
}

type AcceptRideRequest struct {
	DriverID string `json:"driver_id" binding:"required"`
}
