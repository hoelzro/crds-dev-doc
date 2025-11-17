package main

import (
	"net/http"
	"strings"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture status code and bytes written
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(b)
	rr.bytes += n
	return n, err
}

// getClientIP extracts the real client IP from request headers
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("True-Client-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default if WriteHeader is not called
			bytes:          0,
		}

		next.ServeHTTP(recorder, r)

		logger.Debug("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration", time.Since(start),
			"bytes", recorder.bytes,
			"ip", getClientIP(r),
		)
	})
}
