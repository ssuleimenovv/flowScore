package api

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

// CORS lets pages from the allowed origins call the API from another host:
// in production the frontend is on Vercel and the API on Render. Origins are
// host patterns, the same the WebSocket stream checks: "flowscore.vercel.app",
// "flowscore-*.vercel.app", "localhost:5173".
func CORS(next http.Handler, origins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" || !allowed(origin, origins) {
			next.ServeHTTP(w, r)
			return
		}

		h := w.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Add("Vary", "Origin")

		// The browser asks first before a POST with a JSON body (the simulator)
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			h.Set("Access-Control-Allow-Methods", "GET, POST")
			h.Set("Access-Control-Allow-Headers", "Content-Type")
			h.Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// allowed reports whether the origin's host matches one of the patterns.
func allowed(origin string, patterns []string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Host)
	for _, p := range patterns {
		if ok, _ := path.Match(strings.ToLower(p), host); ok {
			return true
		}
	}
	return false
}
