package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/entity"
)

type DB struct {
	db *sql.DB
}

func (d *DB) CreateRequest(r entity.Request) (int, error) {
	query := `
	INSERT INTO public.request (
		contact,
		contact_type,
		message,
		created_at
	)
	VALUES ($1, $2, $3, $4) 
	RETURNING id`

	var id int
	err := d.db.QueryRow(query, r.Contact, r.ContactType, r.Message, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *DB) GetRequests() ([]entity.Request, error) {
	query := `SELECT id, contact, contact_type, message, created_at FROM public.request`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.Request
	for rows.Next() {
		var r entity.Request
		err := rows.Scan(
			&r.ID,
			&r.Contact,
			&r.ContactType,
			&r.Message,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}

func (d *DB) GetPasswordHash(userName string) (*entity.User, error) {
	query := `
	SELECT 
		password_hash 
	FROM 
	public.users 
	WHERE user_name = $1`

	row := d.db.QueryRow(query, userName)
	user := entity.User{}

	err := row.Scan(
		&user.PasswordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("пользователь %v не найден", userName)
		}
		return nil, fmt.Errorf("ошибка при чтении данных пользователя: %w", err)
	}

	return &user, nil
}

func (d *DB) DeleteRequest(id int) error {
	query := `DELETE FROM public.request WHERE id = $1`

	_, err := d.db.Exec(query, id)

	return err
}
