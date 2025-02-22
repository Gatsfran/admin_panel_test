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

type ClientOrder struct {
	ID          int         `json:"id"`
	Contact     string      `json:"contact"`
	ContactType ContactType `json:"contact_type"`
	Message     string      `json:"message"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (r ClientOrder) String() string {
	return fmt.Sprintf(`Тип заявки: %s`, r.ContactType)
}

type Outbox struct {
	OrderID int  `json:"order_id"`
	IsSent  bool `json:"is_sent"`
}
