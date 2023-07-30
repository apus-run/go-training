package dao

import (
	"context"

	"project-layout/internal/repository"
	"project-layout/internal/repository/dao/model"
)

var _ UserDAO = (*userDAO)(nil)

type UserDAO interface {
	Insert(ctx context.Context, user *model.User) (uint64, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userEntity *model.User) error
	Get(ctx context.Context, id uint64) (*model.User, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.User, error)
}

type userDAO struct {
	data *repository.Data
}

func (u *userDAO) Insert(ctx context.Context, userModel *model.User) (uint64, error) {
	err := u.data.DB.WithContext(ctx).Create(userModel).Error
	return userModel.ID, err
}

func (u *userDAO) Update(ctx context.Context, userModel *model.User) error {
	err := u.data.DB.WithContext(ctx).Updates(userModel).Error
	return err
}

func (u *userDAO) Delete(ctx context.Context, userModel *model.User) error {
	return u.data.DB.WithContext(ctx).Delete(userModel).Error
}

func (u *userDAO) Get(ctx context.Context, id uint64) (*model.User, error) {
	userModel := new(model.User)
	err := u.data.DB.WithContext(ctx).Where("id=?", id).First(userModel).Error
	return userModel, err
}

func (u *userDAO) GetByUserID(ctx context.Context, userID string) ([]*model.User, error) {
	users := make([]*model.User, 0)
	err := u.data.DB.WithContext(ctx).Where("user_id=?", userID).Find(users).Error
	return users, err
}
