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
			// Разрешаем запросы от всех доменов, если это не production
			if !cfg.Server.IsProduction {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// В production разрешаем только указанные домены
				for _, origin := range cfg.Server.CORS.Allow_origins {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			// Разрешаем указанные методы
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Server.CORS.Allow_methods, ", "))

			// Разрешаем указанные заголовки
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.Server.CORS.Allow_headers, ", "))

			// Разрешаем передачу cookies и авторизационных заголовков
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Если это предварительный запрос (OPTIONS), завершаем его
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Передаем запрос дальше по цепочке middleware
			next.ServeHTTP(w, r)
		})
	}
}