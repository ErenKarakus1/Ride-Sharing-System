package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var request CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.repository.Create(ctx.Request.Context(), request.ToUser())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (h *Handler) List(ctx *gin.Context) {
	users, err := h.repository.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *Handler) Get(ctx *gin.Context) {
	found, err := h.repository.Get(ctx.Request.Context(), ctx.Param("id"))
	if errors.Is(err, ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, found)
}

func (h *Handler) Update(ctx *gin.Context) {
	var request UpdateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repository.Update(ctx.Request.Context(), ctx.Param("id"), request)
	if errors.Is(err, ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"display_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Role        Role   `json:"role" binding:"required"`
}

func (r CreateUserRequest) ToUser() User {
	return User{
		Email:       strings.TrimSpace(r.Email),
		DisplayName: strings.TrimSpace(r.DisplayName),
		PhoneNumber: strings.TrimSpace(r.PhoneNumber),
		Role:        r.Role,
	}
}

type UpdateUserRequest struct {
	DisplayName *string `json:"display_name"`
	PhoneNumber *string `json:"phone_number"`
}
