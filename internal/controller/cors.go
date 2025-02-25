package controller

import (
	"net/http"
	"strings"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/gorilla/mux"
)

func CORSMiddleware(cfg *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !cfg.Server.IsProduction {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				for _, origin := range cfg.Server.CORS.Allow_origins {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
			
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Server.CORS.Allow_methods, ", "))

			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.Server.CORS.Allow_headers, ", "))

			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
