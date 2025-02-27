package telegram

import (
	"fmt"
	"log"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramBot(cfg *config.Telegram) (*TelegramBot, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("токен бота не может быть пустым")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании бота: %w", err)
	}

	bot.Debug = true
	log.Printf("Авторизован как бот: %s", bot.Self.UserName)

	return &TelegramBot{
		bot: bot,
	}, nil
}

func (t *TelegramBot) SendMessage(chatID int64, text string) error {
	if chatID == 0 {
		return fmt.Errorf("chat_id не может быть пустым")
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}

	return nil
}
