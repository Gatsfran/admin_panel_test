package entity

import (
	"fmt"
	"time"
)

type User struct {
	ID           int    `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

type Request struct {
	ID          int       `json:"id"`
	Contact     string    `json:"contact"`
	ContactType string    `json:"contact_type"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r Request) String() string {
	return fmt.Sprintf(`
	Информация о запросе:
	========================
	ID:             %d
	Контакты:       %s
	Сообщение       %s
	Время создания: %s
	========================
	`, r.ID, r.Contact, r.Message, r.CreatedAt)
}

func (u User) String() string {
	return fmt.Sprintf(`
	Информация о запросе:
	========================
	ID:               %d
	Имя пользователя: %s
	Пороль            %s
	========================
	`, u.ID, u.UserName, u.PasswordHash)
}
