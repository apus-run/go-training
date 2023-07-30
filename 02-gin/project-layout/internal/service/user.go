package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
)

type UserService struct {
	repo *repository.UserRepo
	log  *logrus.Logger
}

func NewUserService(repo *repository.UserRepo, logger *logrus.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  logger,
	}
}

func (svc *UserService) Login(ctx context.Context, username, password string) (*entity.User, error) {

	return &entity.User{}, nil
}

func (svc *UserService) Register(ctx context.Context, userEntity *entity.User) (*entity.User, error) {
	return &entity.User{}, nil
}
