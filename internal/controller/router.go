package controller

import (
	"net/http"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/Gatsfran/admin_panel_test/internal/repo"
	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
	db     *repo.DB
	cfg    *config.Config
}

func New(db *repo.DB, cfg *config.Config) *Router {
	r := mux.NewRouter()
	router := &Router{
		router: r,
		db:     db,
		cfg:    cfg,
	}

	router.registerAuthRoutes()
	router.registerClientOrderRoutes()

	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
