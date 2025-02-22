package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
    bot *tgbotapi.BotAPI
    chatID int64
}

func NewTelegramBot(token string, chatID int64) (*TelegramBot, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, fmt.Errorf("ошибка при создании бота: %w", err)
    }
    return &TelegramBot{bot: bot, chatID: chatID}, nil
}

func (t *TelegramBot) SendMessage(text string) error {
    msg := tgbotapi.NewMessage(t.chatID, text)
    _, err := t.bot.Send(msg)
    if err != nil {
        return fmt.Errorf("ошибка при отправке сообщения: %w", err)
    }
    return nil
}