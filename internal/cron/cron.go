package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/repo"
	"github.com/Gatsfran/admin_panel_test/internal/telegram"
)

type Cron struct {
	outboxRepo *repo.DB
	tgBot      *telegram.TelegramBot
	interval   time.Duration
}

func NewCron(outboxRepo *repo.DB, tgBot *telegram.TelegramBot, interval time.Duration) *Cron {
	return &Cron{
		outboxRepo: outboxRepo,
		tgBot:      tgBot,
		interval:   interval,
	}
}

func (c *Cron) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.processOutbox(ctx)
		case <-ctx.Done():
			log.Println("Cron процесс остановлен")
			return
		}
	}
}

func (c *Cron) processOutbox(ctx context.Context) {
	log.Println("Проверка неотправленных заявок...")

	outboxItems, err := c.outboxRepo.GetUnsentOrders(ctx)
	if err != nil {
		log.Printf("Ошибка при получении заявок: %v", err)
		return
	}

	for _, item := range outboxItems {
		// Отправляем сообщение в Telegram
		message := fmt.Sprintf("Заявка #%d готова к обработке", item.OrderID)
		if err := c.tgBot.SendMessage(message); err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			continue
		}

		if err := c.outboxRepo.MarkAsSent(ctx, item.OrderID); err != nil {
			log.Printf("Ошибка при обновлении заявки: %v", err)
			continue
		}

		log.Printf("Заявка #%d успешно отправлена", item.OrderID)
	}
}