package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/loissascha/go-logger/logger"
)

func corsMiddleware(next http.Handler) http.Handler {
	rawAllowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	allowAllOrigins := false
	if rawAllowedOrigins == "*" {
		allowAllOrigins = true
	}
	if rawAllowedOrigins == "" && os.Getenv("APP_ENV") == "production" {
		logger.Warning(nil, "Allowed origins is not set! Please make sure to configure your .env file!")
	} else if rawAllowedOrigins == "" {
		logger.Warning(nil, "Allowed origins is not set! Allowing development hosts!")
		rawAllowedOrigins = "http://localhost:4321,http://localhost:4322,http://127.0.0.1:4321,http://127.0.0.1:4322,http://localhost:5173,http://127.0.0.1:5173,http://localhost:3000,http://127.0.0.1:3000"
	}
	allowedOriginsList := strings.Split(rawAllowedOrigins, ",")
	allowedOriginsMap := make(map[string]bool)
	for _, origin := range allowedOriginsList {
		trimmedOrigin := strings.TrimSpace(origin)
		if trimmedOrigin != "" {
			allowedOriginsMap[trimmedOrigin] = true
		}
	}
	if len(allowedOriginsMap) == 0 && os.Getenv("APP_ENV") == "production" {
		logger.Fatal(nil, "FATAL: No valid allowed origins configured for CORS in production.")
		panic("CORS not configured!")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestOrigin := r.Header.Get("Origin")
		originAllowed := false
		// logger.Info(nil, "Request Origin: {origin}", requestOrigin)

		if requestOrigin != "" {
			if _, ok := allowedOriginsMap[requestOrigin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
				originAllowed = true
			} else {
				if os.Getenv("APP_ENV") != "production" && (strings.HasPrefix(requestOrigin, "http://localhost:") || strings.HasPrefix(requestOrigin, "http://127.0.0.1:")) {
					w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
					originAllowed = true
					// logger.Info(nil, "CORS: Development: Allowed reflected origin {origin}", requestOrigin)
				}
			}
		} else {
			// No Origin header usually means same-origin or a server-to-server request.
			originAllowed = true
		}

		if allowAllOrigins {
			originAllowed = true
			if requestOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
			}
		}

		if originAllowed || r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH, HEAD")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Origin")
		}

		if r.Method == "OPTIONS" {
			if originAllowed { // Only respond positively to OPTIONS if the origin will be allowed
				w.WriteHeader(http.StatusNoContent)
			} else {
				// Origin not allowed, so preflight should effectively fail.
				http.Error(w, "CORS: Origin not allowed", http.StatusForbidden)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
