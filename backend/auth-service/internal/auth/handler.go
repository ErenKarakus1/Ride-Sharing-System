package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(ctx.Request.Context(), request)
	if errors.Is(err, ErrEmailAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register account"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(ctx.Request.Context(), request)
	if errors.Is(err, ErrInvalidCredentials) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
