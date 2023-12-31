package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"project-layout/internal/infra"
	"project-layout/internal/repository/dao/model"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrUserDuplicate = errors.New("email or phone number already exists")
var DuplicateEntryErrCode uint16 = 1062

// UserDAO ... 数据访问层相关接口定义, 使用 DB 风格的命名
type UserDAO interface {
	Insert(ctx context.Context, userModel model.User) (uint64, error)
	Update(ctx context.Context, userModel model.User) error
	Delete(ctx context.Context, userModel model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	FindByWechat(ctx context.Context, openId string) (*model.User, error)
}

type userDAO struct {
	data *infra.Data
}

func NewUserDAO(data *infra.Data) UserDAO {
	return &userDAO{
		data: data,
	}
}

func (u *userDAO) Insert(ctx context.Context, userModel model.User) (uint64, error) {
	now := time.Now().UnixMilli()
	userModel.CreatedTime = now
	userModel.UpdatedTime = now

	err := u.data.DB.WithContext(ctx).Create(&userModel).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == DuplicateEntryErrCode {
			// 邮箱或者手机号冲突
			return 0, ErrUserDuplicate
		}
	}

	return userModel.ID, err
}

func (u *userDAO) Update(ctx context.Context, userModel model.User) error {
	// 这种写法是很不清晰的，因为它依赖了 gorm 的两个默认语义
	// 会使用 ID 来作为 WHERE 条件
	// 会使用非零值来更新
	// 另外一种做法是显式指定只更新必要的字段，
	// 那么这意味着 DAO 和 service 中非敏感字段语义耦合了
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

func (u *userDAO) FindByWechat(ctx context.Context, openId string) (*model.User, error) {
	var userModel model.User
	err := u.data.DB.WithContext(ctx).First(&userModel, "wechat_open_id = ?", openId).Error
	return &userModel, err
}
