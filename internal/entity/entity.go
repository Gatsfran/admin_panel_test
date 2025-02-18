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

type ContactType string

const (
	Email    ContactType = "email"
	Phone    ContactType = "phone"
	Telegram ContactType = "telegram"
)

type Order struct {
	ID          int         `json:"id"`
	Contact     string      `json:"contact"`
	ContactType ContactType `json:"contact_type"`
	Message     string      `json:"message"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (r Order) String() string {
	return fmt.Sprintf(`Тип заявки: %s`, r.ContactType)
}
