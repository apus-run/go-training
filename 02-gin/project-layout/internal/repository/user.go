package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository/cache"
	"project-layout/internal/repository/dao"
	"project-layout/internal/repository/dao/model"
	"project-layout/pkg/log"
)

var ErrUserDuplicateEmailOrPhone = dao.ErrUserDuplicateEmailOrPhone
var ErrUserDataNotFound = dao.ErrRecordNotFound

type UserRepository interface {
	// Save 是 Create and Update 统称
	Save(ctx context.Context, userEntity entity.User) error

	SaveAndCache(ctx context.Context, userEntity entity.User) error

	Remove(ctx context.Context, userEntity entity.User) error
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	FindUserPage(ctx context.Context, name string, page int64, size int64) ([]*entity.User, uint, bool, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// userRepository 使用了缓存
type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache

	log *log.Logger
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache, logger *log.Logger) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
		log:   logger,
	}
}

// Save 保存用户信息到数据库
func (ur *userRepository) Save(ctx context.Context, userEntity entity.User) error {
	ur.log.Info("create or update user")
	// Map the data from Entity to DO
	userModel := model.User{}
	userModel, _ = userModel.FromEntity(userEntity).(model.User)

	if userModel.ID == 0 {
		// Insert
		_, err := ur.dao.Insert(ctx, userModel)
		if err != nil {
			return err
		}
	} else {
		// Update
		return ur.dao.Update(ctx, userModel)
	}

	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()
	userEntity = newEntity

	return nil
}

// SaveAndCache 保存用户信息到数据库和缓存
func (ur *userRepository) SaveAndCache(ctx context.Context, userEntity entity.User) error {
	ur.log.Info("save and cache user")
	// Map the data from Entity to DO
	userModel := model.User{}
	userModel, _ = userModel.FromEntity(userEntity).(model.User)

	// Save the data into DB
	id, err := ur.dao.Insert(ctx, userModel)
	if err != nil {
		return err
	}

	userModel.ID = id
	err = ur.cache.SetOjb(ctx, userEntity)
	if err != nil {
		return err
	}

	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()
	userEntity = newEntity

	return nil
}

// Remove 从数据库和缓存中删除用户信息
func (ur *userRepository) Remove(ctx context.Context, user entity.User) error {
	ur.log.Info("remove user")
	// Map the data from Entity to DO
	userModel := model.User{}
	userModel, _ = userModel.FromEntity(user).(model.User)

	// Remove the data from DB
	err := ur.dao.Delete(ctx, userModel)
	if err != nil {
		return err
	}

	// Remove the data from cache
	err = ur.cache.Del(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// FindByID 从数据库和缓存中获取用户信息
func (ur *userRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	ur.log.Info("find user")
	userEntity := &entity.User{}

	// 1. 先从缓存中获取
	res, err := ur.cache.Get(ctx, id)
	if err == nil {
		return &res, err
	}

	// 2. 缓存中没有，从数据库中获取
	userModel, err := ur.dao.FindByID(ctx, id)

	// 对错误进行包装, 尽量使用通一的错误处理(Sentinel error), 屏蔽掉不同数据库的报错差异性
	// 这里不只是返回这一种错误, 在上层打日志的时候还要看到底层出的是那种错误, 所以原始的错误信息也需要保留.
	if err != nil {
		if errors.Is(err, ErrUserDataNotFound) {
			return userEntity, errors.Wrap(ErrUserDataNotFound, fmt.Sprintf("此用户不存在, err: %v", err))
		} else {
			return userEntity, err
		}
	}

	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()

	// 3. 更新缓存
	go func() {
		err = ur.cache.SetOjb(ctx, newEntity)
		if err != nil {
			// 打日志, 做监控, 可以推断出缓存服务是否正常
			ur.log.Error("set cache failed", zap.Error(err))
		}
	}()

	return &newEntity, nil
}

func (ur *userRepository) FindByIdV1(ctx context.Context, id uint64) (*entity.User, error) {
	u, err := ur.cache.Get(ctx, id)
	switch err {
	case nil:
		return &u, err
	case cache.ErrKeyNotExist:
		userModel, err := ur.dao.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		// Map fresh record's data into Entity
		newEntity := userModel.ToEntity()
		if err != nil {
			return nil, err
		}
		return &newEntity, nil
	default:
		return &entity.User{}, err
	}
}

// FindUserPage 分页查询用户信息
func (ur *userRepository) FindUserPage(ctx context.Context,
	name string, page int64, size int64) ([]*entity.User, uint, bool, error) {
	ur.log.Info("find user page by user name")
	return nil, 0, true, nil
}

// FindByPhone 根据手机号查询用户信息
func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	userModel, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return &entity.User{}, err
	}
	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()
	if err != nil {
		return &entity.User{}, err
	}

	return &newEntity, nil
}

// FindByEmail 根据邮箱查询用户信息
func (ur *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return &entity.User{}, err
	}
	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()
	if err != nil {
		return &entity.User{}, err
	}

	return &newEntity, nil
}
