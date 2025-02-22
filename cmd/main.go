package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/Gatsfran/admin_panel_test/internal/controller"
	"github.com/Gatsfran/admin_panel_test/internal/repo"
)

func main() {
	cfg := config.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repo.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	router := controller.New(db, cfg)

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
