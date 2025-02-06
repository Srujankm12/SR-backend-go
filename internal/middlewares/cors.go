package middlewares

import (
	"net/http"
)

// CorsMiddleware sets CORS headers to allow all origins, methods, and headers
func CorsMiddleware(ah http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers (or use wildcard '*' for all headers)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Allow credentials (optional, if your frontend sends cookies or authentication tokens)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests (OPTIONS requests)
		if r.Method == "OPTIONS" {
			// Respond to preflight requests immediately with a 204 No Content status
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ah.ServeHTTP(w, r)
	})
}
