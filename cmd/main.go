package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/Gatsfran/admin_panel_test/internal/controller"
	"github.com/Gatsfran/admin_panel_test/internal/cron"
	"github.com/Gatsfran/admin_panel_test/internal/repo"
	"github.com/Gatsfran/admin_panel_test/internal/telegram"
)

func main() {
	cfg := config.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repo.New(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	router := controller.New(db, cfg)

	tgBot, err := telegram.NewTelegramBot(cfg.Telegram)
	if err != nil {
		log.Fatalf("Ошибка при создании Telegram бота: %v", err)
	}

	cronProcess := cron.NewCron(db, tgBot, cfg.Telegram.ChatID, 1*time.Minute)
	go cronProcess.Start(ctx)

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
