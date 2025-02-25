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
	chatID     int64
	interval   time.Duration
}

func NewCron(outboxRepo *repo.DB, tgBot *telegram.TelegramBot, chatID int64, interval time.Duration) *Cron {
	return &Cron{
		outboxRepo: outboxRepo,
		tgBot:      tgBot,
		chatID:     chatID,
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
	log.Printf("Проверка неотправленных заявок...")
	outboxItems, err := c.outboxRepo.GetUnsentOrders(ctx)
	if err != nil {
		log.Printf("Ошибка при получении заявок: %v", err)
		return
	}

	for _, item := range outboxItems {
		message := fmt.Sprintf("Поступила заявка:%s", item)
		if err := c.tgBot.SendMessage(c.chatID, message); err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			continue
		}
		if err := c.outboxRepo.MarkAsSent(ctx, item.ID); err != nil {
			log.Printf("Ошибка при обновлении заявки: %v", err)
			continue
		}
		log.Printf("Заявка #%d успешно отправлена", item.ID)
	}
}
