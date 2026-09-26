package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httprequest"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httpresponse"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/dto"
)

type userService interface {
	GetUsers(ctx context.Context) ([]dto.UserResponse, error)
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.UserResponse, error)
}

type Handler struct {
	service userService
}

func New(service userService) *Handler {
	return &Handler{service: service}
}

const (
	maxRequestBodySizeMB = 1
	maxRequestBodySize   = maxRequestBodySizeMB * 1_000_000
	// maxRequestBodySize limits incoming request bodies to 1 MB (1,000,000 bytes).
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		slog.Error("Failed to fetch users", "action", "USERS_FETCH_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "USERS_FETCH_FAILED", "Failed to fetch users")
		return
	}

	httpresponse.Success(w, http.StatusOK, "Users fetched successfully", users)
	slog.Info("Users fetched successfully", "action", "USERS_FETCH_SUCCESS", "count", len(users))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	// 1. Parse + validate JSON structure
	var request dto.CreateUserRequest
	isJSONValid := httprequest.ParseAndValidateJSON(w, r, &request, maxRequestBodySize)
	if !isJSONValid {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request")
		return
	}

	// 2. Validate business/input fields
	if err := request.Validate(); err != nil {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			err.Error(),
		)
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), request)
	if err != nil {
		slog.Error("Failed to create user", "action", "USER_CREATE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_CREATE_FAILED", "Failed to create user")
		return
	}

	httpresponse.Success(w, http.StatusCreated, "User created successfully", createdUser)
}
