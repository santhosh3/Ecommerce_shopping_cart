package globalmiddlewares

import (
	"net/http"
)


// CORS middleware
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allowed methods
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allowed headers

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
            return
        }
        next.ServeHTTP(w, r) // Call the next handler
    })
}
