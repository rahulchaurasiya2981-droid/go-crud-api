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

func logHandlerStarted(r *http.Request) time.Time {
	start := time.Now()
	slog.Info("API handler started",
		"action", "API_HANDLER_START",
		"method", r.Method,
		"path", r.URL.Path,
	)
	return start
}

func logHandlerCompleted(start time.Time, r *http.Request) {
	slog.Info("API handler completed",
		"action", "API_HANDLER_COMPLETED",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func logHandlerFailed(start time.Time, r *http.Request, err error) {
	slog.Warn("API handler failed",
		"action", "API_HANDLER_FAILED",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_ms", time.Since(start).Milliseconds(),
		"error", err,
	)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	start := logHandlerStarted(r)

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		slog.Error("Failed to fetch users", "action", "USERS_FETCH_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusInternalServerError, "USERS_FETCH_FAILED", "Failed to fetch users")
		return
	}

	slog.Info("Users fetched successfully", "action", "USERS_FETCH_SUCCESS", "count", len(users))
	logHandlerCompleted(start, r)
	httpresponse.Success(w, http.StatusOK, "Users fetched successfully", users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := logHandlerStarted(r)

	var request dto.CreateUserRequest
	if !httprequest.ParseAndValidateJSON(w, r, &request, maxRequestBodySize) {
		err := errors.New("invalid JSON request")
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request")
		return
	}

	if err := request.Validate(); err != nil {
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), request)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			slog.Warn("User creation conflict",
				"action", "USER_CREATE_CONFLICT",
				"method", r.Method,
				"path", r.URL.Path,
				"error", err,
			)
			logHandlerFailed(start, r, err)
			httpresponse.Error(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already exists")
			return
		}

		slog.Error("Failed to create user",
			"action", "USER_CREATE_ERROR",
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		)
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_CREATE_FAILED", "Failed to create user")
		return
	}
	logHandlerCompleted(start, r)
	httpresponse.Success(w, http.StatusOK, "User created successfully", createdUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	start := logHandlerStarted(r)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		if err == nil {
			err = errors.New("user ID must be a positive integer")
		}
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a positive integer")
		return
	}

	deletedUser, err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			logHandlerFailed(start, r, err)
			httpresponse.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}

		slog.Error("Failed to delete user", "action", "USER_DELETE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_DELETE_FAILED", "Failed to delete user")
		return
	}
	logHandlerCompleted(start, r)
	httpresponse.Success(w, http.StatusOK, "User deleted successfully", deletedUser)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	start := logHandlerStarted(r)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		if err == nil {
			err = errors.New("user ID must be a positive integer")
		}
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID must be a positive integer")
		return
	}

	var request dto.UpdateUserRequest
	if !httprequest.ParseAndValidateJSON(w, r, &request, maxRequestBodySize) {
		err := errors.New("invalid JSON request")
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request")
		return
	}

	if err := request.Validate(); err != nil {
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), id, request)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			logHandlerFailed(start, r, err)
			httpresponse.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
			return
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			logHandlerFailed(start, r, err)
			httpresponse.Error(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already exists")
			return
		}

		slog.Error("Failed to update user", "action", "USER_UPDATE_ERROR", "method", r.Method, "path", r.URL.Path, "error", err)
		logHandlerFailed(start, r, err)
		httpresponse.Error(w, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Failed to update user")
		return
	}

	logHandlerCompleted(start, r)
	httpresponse.Success(w, http.StatusOK, "User updated successfully", updatedUser)
}
