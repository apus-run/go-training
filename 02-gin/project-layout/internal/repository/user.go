package repository

import (
	"context"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository/cache"
	"project-layout/internal/repository/dao"
	"project-layout/internal/repository/dao/model"
	"project-layout/pkg/log"
)

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

	// 1. 先从缓存中获取
	res, err := ur.cache.Get(ctx, id)
	if err == nil {
		return &res, err
	}

	// 2. 缓存中没有，从数据库中获取
	userModel, err := ur.dao.FindByID(ctx, id)
	if err != nil {
		return &entity.User{}, err
	}

	// Map fresh record's data into Entity
	newEntity := userModel.ToEntity()
	if err != nil {
		return &entity.User{}, err
	}

	// 3. 更新缓存
	_ = ur.cache.SetOjb(ctx, newEntity)

	return &newEntity, nil
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
