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

func New(cfg *config.Config) (*DB, error) {

	pool, err := pgxpool.New(context.Background(), cfg.GetPostgresConnectionString())
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пула соединений: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка при проверке соединения с базой данных: %w", err)
	}

	return &DB{pool: pool}, nil
}
func (d *DB) CreateOrder(ctx context.Context, r *entity.Order) error {
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
	err := d.pool.QueryRow(ctx, query, r.Contact, r.ContactType, r.Message, time.Now()).Scan(&id)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	return nil
}

func (d *DB) ListOrder(ctx context.Context) ([]entity.Order, error) {
	query := `SELECT id, contact, contact_type, message, created_at FROM public.request`

	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения запросов: %w", err)
	}
	defer rows.Close()

	var requests []entity.Order
	for rows.Next() {
		var r entity.Order
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
		requests = append(requests, r)
	}
	return requests, nil
}

func (d *DB) GetPasswordHash(ctx context.Context, userName string) error {
	query := `
	SELECT 
		password_hash 
	FROM public.users 
	WHERE user_name = $1`

	row := d.pool.QueryRow(ctx, query, userName)
	user := entity.User{}

	err := row.Scan(
		&user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("ошибка при получении хэша пароля: %w", err)
	}
	return nil
}

func (d *DB) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM public.request WHERE id = $1`

	_, err := d.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении заявки: %w", err)
	}
	return nil
}
