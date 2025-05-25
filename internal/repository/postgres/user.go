package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"raspyxTelegramNotifier/internal/domain/models"
	"raspyxTelegramNotifier/internal/repository"
	"strings"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	const op = "repository.postgres.UserRepository.Create"

	query := `INSERT INTO users (telegram_id)
			  VALUES ($1)`
	_, err := r.db.Exec(ctx, query, user.TelegramID)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return fmt.Errorf("%s: %w", op, repository.ErrExist)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepository) Get(ctx context.Context) ([]models.User, error) {
	const op = "repository.postgres.UserRepository.Get"

	query := `SELECT telegram_id
			  FROM users`
	rows, err := r.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.TelegramID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, tid int64) error {
	const op = "repository.postgres.UserRepository.Delete"

	query := `DELETE FROM users
			  WHERE telegram_id = $1`
	result, err := r.db.Exec(ctx, query, tid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNotFound)
	}

	return nil
}
