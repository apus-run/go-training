package service

import (
	"context"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
	"project-layout/pkg/log"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
	Register(ctx context.Context, userEntity *entity.User) (*entity.User, error)
}

type userService struct {
	repo repository.UserRepo
	log  *log.Logger
}

func NewUserService(repo repository.UserRepo, logger *log.Logger) UserService {
	return &userService{
		repo: repo,
		log:  logger,
	}
}

func (svc *userService) Login(ctx context.Context, username, password string) (*entity.User, error) {

	return &entity.User{}, nil
}

func (svc *userService) Register(ctx context.Context, userEntity *entity.User) (*entity.User, error) {
	return &entity.User{}, nil
}
