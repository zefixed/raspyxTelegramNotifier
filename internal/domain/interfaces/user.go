package interfaces

import (
	"context"
	"raspyxTelegramNotifier/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Get(ctx context.Context) ([]models.User, error)
	//GetByTelegramID(ctx context.Context, tid int64) (*models.User, error)
	//Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, tid int64) error
}
