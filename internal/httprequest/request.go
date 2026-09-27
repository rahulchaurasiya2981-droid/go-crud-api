package httprequest

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/logger"
)

var errTrailingJSONData = errors.New("request body must contain one JSON value")

func ParseAndValidateJSON(w http.ResponseWriter, r *http.Request, destination any, maxBodySize int64) bool {
	start := time.Now()
	slog.Info(logger.MsgHTTPRequestJSONValidationStart,
		"action", logger.ActionHTTPRequestJSONValidationStart,
		"method", r.Method,
		"path", r.URL.Path,
	)

	defer func() {
		slog.Info(logger.MsgHTTPRequestJSONValidationEnd,
			"action", logger.ActionHTTPRequestJSONValidationEnd,
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}()

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// 1. Initial Decode
	if err := decoder.Decode(destination); err != nil {
		if isRequestBodyTooLarge(err) {
			slog.Warn(logger.MsgHTTPRequestBodyTooLarge,
				"action", logger.ActionHTTPRequestBodyTooLarge,
				"method", r.Method,
				"path", r.URL.Path,
				"error", err,
			)
			return false
		}

		slog.Error(logger.MsgHTTPRequestJSONDecodeError,
			"action", logger.ActionHTTPRequestJSONDecodeError,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		)
		return false
	}

	// 2. Trailing Data Check
	var trailingData any
	if err := decoder.Decode(&trailingData); err != io.EOF {
		if isRequestBodyTooLarge(err) {
			slog.Warn(logger.MsgHTTPRequestBodyTooLarge,
				"action", logger.ActionHTTPRequestBodyTooLarge,
				"method", r.Method,
				"path", r.URL.Path,
				"error", err,
			)
			return false
		}

		if err == nil {
			err = errTrailingJSONData
		}

		slog.Error(logger.MsgHTTPRequestTrailingData,
			"action", logger.ActionHTTPRequestTrailingData,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		)
		return false
	}

	return true
}

func isRequestBodyTooLarge(err error) bool {
	var maxBytesError *http.MaxBytesError
	return errors.As(err, &maxBytesError)
}
