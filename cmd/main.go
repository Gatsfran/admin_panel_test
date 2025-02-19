package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/Gatsfran/admin_panel_test/internal/controller"
	"github.com/Gatsfran/admin_panel_test/internal/repo"
)

func main() {
	cfg := &config.Postgres{
		Host:     "localhost",
		Port:     "15423",
		Username: "postgres",
		Password: "your_password",
		Database: "admin_db",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repo.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	router := controller.New(db)
	router.RegisterAuthRoutes() // Регистрация ручки /login

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}