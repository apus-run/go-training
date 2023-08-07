package dao

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"project-layout/internal/repository"
	"project-layout/internal/repository/dao/model"
)

var ErrRecordNotFound = errors.New("record not found")

// UserDAO ... 数据访问层相关接口定义, 使用 DB 风格的命名
type UserDAO interface {
	Insert(ctx context.Context, userModel model.User) (uint64, error)
	Update(ctx context.Context, userModel model.User) error
	Delete(ctx context.Context, userModel model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type userDAO struct {
	data *repository.Data
}

func NewUserDAO(data *repository.Data) UserDAO {
	return &userDAO{
		data: data,
	}
}

func (u *userDAO) Insert(ctx context.Context, userModel model.User) (uint64, error) {
	now := time.Now().UnixMilli()
	userModel.CreatedTime = now
	userModel.UpdatedTime = now

	err := u.data.DB.WithContext(ctx).Create(&userModel).Error
	return userModel.ID, err
}

func (u *userDAO) Update(ctx context.Context, userModel model.User) error {
	err := u.data.DB.WithContext(ctx).Updates(&userModel).Error
	return err
}

func (u *userDAO) Delete(ctx context.Context, userModel model.User) error {
	return u.data.DB.WithContext(ctx).Delete(&userModel, "id = ?", userModel.ID).Error
}

func (u *userDAO) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var userModel model.User
	err := u.data.DB.WithContext(ctx).First(&userModel, "id = ?", id).Error
	return &userModel, err
}

func (u *userDAO) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var userModel model.User
	err := u.data.DB.WithContext(ctx).First(&userModel, "phone = ?", phone).Error
	return &userModel, err
}

func (u *userDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userModel model.User
	err := u.data.DB.WithContext(ctx).First(&userModel, "email = ?", email).Error
	return &userModel, err
}
