package payment

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

func (h *Handler) Authorize(ctx *gin.Context) {
	var request AuthorizeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.Authorize(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize payment"})
		return
	}

	ctx.JSON(http.StatusCreated, payment)
}

func (h *Handler) Get(ctx *gin.Context) {
	payment, err := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	h.writePaymentResult(ctx, payment, err)
}

func (h *Handler) Capture(ctx *gin.Context) {
	payment, err := h.service.Capture(ctx.Request.Context(), ctx.Param("id"))
	h.writePaymentResult(ctx, payment, err)
}

func (h *Handler) Refund(ctx *gin.Context) {
	payment, err := h.service.Refund(ctx.Request.Context(), ctx.Param("id"))
	h.writePaymentResult(ctx, payment, err)
}

func (h *Handler) writePaymentResult(ctx *gin.Context, payment Payment, err error) {
	if errors.Is(err, ErrPaymentNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	if errors.Is(err, ErrInvalidPaymentTransition) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "invalid payment status transition"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

type AuthorizeRequest struct {
	RideID   string  `json:"ride_id" binding:"required"`
	RiderID  string  `json:"rider_id" binding:"required"`
	DriverID *string `json:"driver_id"`
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency"`
}
