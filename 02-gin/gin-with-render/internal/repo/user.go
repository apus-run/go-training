package repo

import (
	"context"

	"gin-with-render/internal/domain"
)

type UserRepo interface {
	Save(ctx context.Context, userDomain domain.User) error
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (u userRepo) Save(ctx context.Context, userDomain domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}
