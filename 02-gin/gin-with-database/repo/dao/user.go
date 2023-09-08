package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"gin-with-database/repo/dao/model"
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
}

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{
		db: db,
	}
}

func (ud *userDAO) Insert(ctx context.Context, userModel model.User) (uint64, error) {
	now := time.Now().UnixMilli()
	userModel.CreatedTime = now
	userModel.UpdatedTime = now

	err := ud.db.WithContext(ctx).Create(&userModel).Error

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == DuplicateEntryErrCode {
			// 邮箱或者手机号冲突
			return 0, ErrUserDuplicate
		}
	}

	return userModel.ID, err
}

func (ud *userDAO) Update(ctx context.Context, userModel model.User) error {
	// 这种写法是很不清晰的，因为它依赖了 gorm 的两个默认语义
	// 会使用 ID 来作为 WHERE 条件
	// 会使用非零值来更新
	// 另外一种做法是显式指定只更新必要的字段，
	// 那么这意味着 DAO 和 service 中非敏感字段语义耦合了
	err := ud.db.WithContext(ctx).Updates(&userModel).Error
	return err
}

func (ud *userDAO) Delete(ctx context.Context, userModel model.User) error {
	return ud.db.WithContext(ctx).Delete(&userModel, "id = ?", userModel.ID).Error
}

func (ud *userDAO) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var userModel model.User
	err := ud.db.WithContext(ctx).First(&userModel, "id = ?", id).Error
	return &userModel, err
}

func (ud *userDAO) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var userModel model.User
	err := ud.db.WithContext(ctx).First(&userModel, "phone = ?", phone).Error
	return &userModel, err
}

func (ud *userDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userModel model.User
	err := ud.db.WithContext(ctx).First(&userModel, "email = ?", email).Error
	return &userModel, err
}
