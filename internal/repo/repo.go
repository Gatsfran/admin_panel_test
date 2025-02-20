package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Gatsfran/admin_panel_test/internal/config"
	"github.com/Gatsfran/admin_panel_test/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func (d *DB) Close() {
	panic("unimplemented")
}

func New(ctx context.Context, cfg *config.Postgres) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.GetPostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пула соединений: %w", err)
	}

	ctxNew, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := pool.Ping(ctxNew); err != nil {
		return nil, fmt.Errorf("ошибка при проверке соединения с базой данных: %w", err)
	}

	return &DB{pool: pool}, nil
}
func (d *DB) CreateClientOrder(ctx context.Context, r *entity.ClientOrder) error {
	query := `
	INSERT INTO public.client_order (
		contact,
		contact_type,
		message,
		created_at
	)
	VALUES ($1, $2, $3, $4) 
	RETURNING id`

	err := d.pool.QueryRow(ctx, query, r.Contact, r.ContactType, r.Message, time.Now()).Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	return nil
}

func (d *DB) ListClientOrder(ctx context.Context) ([]entity.ClientOrder, error) {
	query := `SELECT id, contact, contact_type, message, created_at FROM public.client_order`

	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения запросов: %w", err)
	}
	defer rows.Close()

	var clientorder []entity.ClientOrder
	for rows.Next() {
		var r entity.ClientOrder
		err := rows.Scan(
			&r.ID,
			&r.Contact,
			&r.ContactType,
			&r.Message,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строк: %w", err)
		}
		clientorder = append(clientorder, r)
	}
	return clientorder, nil
}

func (d *DB) GetPasswordHash(ctx context.Context, userName string) (string, error) {
	query := `
	SELECT 
		password_hash 
	FROM public.users 
	WHERE user_name = $1`

	row := d.pool.QueryRow(ctx, query, userName)
	var passwordHash string

	err := row.Scan(
		passwordHash,
	)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении хэша пароля: %w", err)
	}
	return passwordHash, nil
}

func (d *DB) DeleteClientOrder(ctx context.Context, id int) error {
	query := `DELETE FROM public.client_order WHERE id = $1`

	_, err := d.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении заявки: %w", err)
	}
	return nil
}
