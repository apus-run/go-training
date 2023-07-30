package repository

import (
	"context"

	"github.com/sirupsen/logrus"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository/cache"
	"project-layout/internal/repository/dao"
	"project-layout/internal/repository/dao/model"
)

var _ UserRepo = (*userRepo)(nil)

type UserRepo interface {
	// Save ...
	// Save 是 Create and Update 统称
	Save(ctx context.Context, userEntity *entity.User) error

	SaveAndCache(ctx context.Context, userEntity *entity.User) error

	// Remove ...
	Remove(ctx context.Context, userEntity *entity.User) error

	// Find ...
	Find(ctx context.Context, id uint64) (*entity.User, error)
	// FindByUserID ...
	FindByUserID(ctx context.Context, userID string) ([]*entity.User, error)
	// FindUserPage ...
	FindUserPage(ctx context.Context, id uint, page int64, size int64) (*entity.User, uint, bool, error)
	// FindUserPageByUserId ...
	FindUserPageByUserId(ctx context.Context, userID string, page int64, size int64) (*entity.User, uint, bool, error)
}

type userRepo struct {
	dao   dao.UserDAO
	cache cache.UserCache

	log *logrus.Logger
}

func NewUserRepo(dao dao.UserDAO, cache cache.UserCache, logger *logrus.Logger) UserRepo {
	return &userRepo{
		dao:   dao,
		cache: cache,
		log:   logger,
	}
}

func (ur *userRepo) Save(ctx context.Context, userEntity *entity.User) error {
	ur.log.Info("save user")
	// Map the data from Entity to DO
	userModel := new(model.User)
	err := userModel.FromEntity(userEntity)
	if err != nil {
		return err
	}

	_, err = ur.dao.Insert(ctx, userModel)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) SaveAndCache(ctx context.Context, userEntity *entity.User) error {
	ur.log.Info("save and cache user")
	// Map the data from Entity to DO
	userModel := new(model.User)
	err := userModel.FromEntity(userEntity)
	if err != nil {
		return err
	}

	// Save the data into DB
	id, err := ur.dao.Insert(ctx, userModel)
	if err != nil {
		return err
	}

	userModel.ID = id
	err = ur.cache.Set(ctx, *userEntity)
	if err != nil {
		return err
	}

	// Map fresh record's data into Entity
	newEntity, err := userModel.ToEntity()
	if err != nil {
		return err
	}
	*userEntity = *newEntity
	return nil
}

func (ur *userRepo) Remove(ctx context.Context, user *entity.User) error {
	ur.log.Info("remove user")
	return nil
}

func (ur *userRepo) Find(ctx context.Context, id uint64) (*entity.User, error) {
	ur.log.Info("find user")

	res, err := ur.cache.Get(ctx, id)
	if err == nil {
		return &res, err
	}

	userModel, err := ur.dao.Get(ctx, id)
	if err != nil {
		return &entity.User{}, err
	}

	// Map fresh record's data into Entity
	newEntity, err := userModel.ToEntity()
	if err != nil {
		return &entity.User{}, err
	}

	err = ur.cache.Set(ctx, *newEntity)
	if err != nil {
		return nil, err
	}

	return newEntity, nil
}

func (ur *userRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.User, error) {
	ur.log.Info("find user by user id")
	return nil, nil
}

func (ur *userRepo) FindUserPage(ctx context.Context, id uint, page int64, size int64) (*entity.User, uint, bool, error) {
	ur.log.Info("find user page")

	return nil, 0, true, nil
}

func (ur *userRepo) FindUserPageByUserId(ctx context.Context, userID string, page int64, size int64) (*entity.User, uint, bool, error) {
	ur.log.Info("find user page by user id")
	return nil, 0, true, nil
}
