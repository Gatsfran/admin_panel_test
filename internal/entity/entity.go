package entity

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
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
	Contact     string      `json:"contact" validate:"required"`
	ContactType ContactType `json:"contact_type"`
	Message     string      `json:"message" validate:"required,max=500"`
	CreatedAt   time.Time   `json:"created_at"`
	IsSent      bool        `json:"-"`
}

func (c *ClientOrder) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientOrder) SetContactType() error {
	const (
		emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		phoneRegex    = `^\+?[0-9]{10,15}$`
		telegramRegex = `^@[a-zA-Z0-9_]{5,32}$`
	)

	if regexp.MustCompile(emailRegex).MatchString(c.Contact) {
		c.ContactType = Email
		return nil
	} else if regexp.MustCompile(phoneRegex).MatchString(c.Contact) {
		c.ContactType = Phone
		return nil
	} else if regexp.MustCompile(telegramRegex).MatchString(c.Contact) {
		c.ContactType = Telegram
		return nil
	}

	return fmt.Errorf("неверный формат контакта: %s", c.Contact)
}

func (c ClientOrder) String() string {
	return fmt.Sprintf(`
	========================
	ID заявки: %d
	%s: %s
	Сообщение: %s
	Дата создания: %s
	========================
	`, c.ID, c.ContactType, c.Contact, c.Message, c.CreatedAt)
}
