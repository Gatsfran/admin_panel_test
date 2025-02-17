package main

import (
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID           int
	UserName     string
	PasswordHash string
}

type Request struct {
	ID          int
	Contact     string
	ContactType string
	Message     string
	CreatedAt   time.Time
}


func CreateRequest(req Request) (int, error) {
	query := `INSERT INTO public.request (contact, contact_type, message, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := db.QueryRow(query, req.Contact, req.ContactType, req.Message, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(`INSERT INTO public.outbox (request_id, is_sent) VALUES ($1, false)`, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func DeleteRequest(id int) error {
	query := `DELETE FROM public.request WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func GetRequests() ([]Request, error) {
	query := `SELECT id, contact, contact_type, message, created_at FROM public.request`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var r Request
		err := rows.Scan(&r.ID, &r.Contact, &r.ContactType, &r.Message, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}

func GetPasswordHash(userName string) (string, error) {
	query := `SELECT password_hash FROM public.users WHERE user_name = $1`
	var passwordHash string
	err := db.QueryRow(query, userName).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}