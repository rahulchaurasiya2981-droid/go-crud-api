package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httprequest"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httpresponse"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/dto"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/repository"
)

type userService interface {
	GetUsers(ctx context.Context) ([]dto.UserResponse, error)
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) (dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, request dto.UpdateUserRequest) (dto.UserResponse, error)
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
	slog.Info("Users fetched successfully", "action", "USERS_FETCH_SUCCESS", "count", len(users))
	httpresponse.Success(w, http.StatusOK, "Users fetched successfully", users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	// 1. Parse and validate JSON structure.
	var request dto.CreateUserRequest
	isJSONValid := httprequest.ParseAndValidateJSON(w, r, &request, maxRequestBodySize)
	if !isJSONValid {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request")
		return
	}

	// 2. Validate request/business fields.
	if err := request.Validate(); err != nil {
		httpresponse.Error(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			err.Error(),
		)
		return
	}

	// 3. Create user.
	createdUser, err := h.service.CreateUser(r.Context(), request)
	if err != nil {
		// Duplicate email.
		if errors.Is(err, repository.ErrDuplicateEmail) {
			slog.Warn(
				"User creation conflict",
				"action", "USER_CREATE_CONFLICT",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err,
			)

			httpresponse.Error(
				w,
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"Email already exists",
			)
			return
		}

		// Unexpected server/database error.
		slog.Error(
			"Failed to create user",
			"action", "USER_CREATE_ERROR",
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		)

		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"USER_CREATE_FAILED",
			"Failed to create user",
		)
		return
	}

	// 4. Success.
	httpresponse.Success(w, http.StatusOK, "User created successfully", createdUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a positive integer")
		return
	}

	deletedUser, err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			httpresponse.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}

		slog.Error("Failed to delete user", "action", "USER_DELETE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(
			w,
			http.StatusInternalServerError,
			"USER_DELETE_FAILED",
			"Failed to delete user",
		)
		return
	}

	httpresponse.Success(w, http.StatusOK, "User deleted successfully", deletedUser)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("API handler started", "action", "API_HANDLER_START", "method", r.Method, "path", r.URL.Path)
	defer func() {
		slog.Info("API handler completed", "action", "API_HANDLER_END", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	}()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a positive integer")
		return
	}

	var request dto.UpdateUserRequest
	if !httprequest.ParseAndValidateJSON(w, r, &request, maxRequestBodySize) {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request")
		return
	}
	if err := request.Validate(); err != nil {
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), id, request)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			httpresponse.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			httpresponse.Error(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already exists")
			return
		}

		slog.Error("Failed to update user", "action", "USER_UPDATE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Failed to update user")
		return
	}

	httpresponse.Success(w, http.StatusOK, "User updated successfully", updatedUser)
}
