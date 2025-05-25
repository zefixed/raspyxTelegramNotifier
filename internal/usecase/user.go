package usecase

import (
	"context"
	"fmt"
	"raspyxTelegramNotifier/internal/domain/interfaces"
	"raspyxTelegramNotifier/internal/domain/models"
	"raspyxTelegramNotifier/internal/dto"
)

type UserUseCase struct {
	repo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Create(ctx context.Context, userDTO *dto.CreateUser) error {
	const op = "usecase.user.Create"

	// DTO to model
	user := &models.User{TelegramID: userDTO.TelegramID}

	// Adding user to db
	err := uc.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *UserUseCase) Get(ctx context.Context) ([]models.User, error) {
	const op = "usecase.user.Get"

	// Getting all users from db
	users, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, tid int64) error {
	const op = "usecase.user.Delete"

	// Deleting users from db with given uuid
	err := uc.repo.Delete(ctx, tid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
