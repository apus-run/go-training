package service

import (
	"context"
	"errors"
	"fmt"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
	"project-layout/pkg/log"
)

var ErrUserDuplicate = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

type UserService interface {
	Login(ctx context.Context, email, password string) (*entity.User, error)
	Register(ctx context.Context, user entity.User) (*entity.User, error)
	FindOrCreate(ctx context.Context, phone string) (*entity.User, error)
	Profile(ctx context.Context, id uint64) (*entity.User, error)
	UpdateProfile(ctx context.Context, user entity.User) error
}

type userService struct {
	repo repository.UserRepository

	log *log.Logger
}

func NewUserService(repo repository.UserRepository, logger *log.Logger) UserService {
	return &userService{
		repo: repo,
		log:  logger,
	}
}

func (us *userService) Login(ctx context.Context, email, password string) (*entity.User, error) {
	userEntity := &entity.User{}
	user, err := us.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserDataNotFound) {
		return userEntity, fmt.Errorf("通过邮箱查找用户失败: %v", err)
	}
	verify := userEntity.VerifyPassword(user.Password(), password)
	if !verify {
		return userEntity, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (us *userService) Register(ctx context.Context, user entity.User) (*entity.User, error) {
	userEntity := &entity.User{}
	hash, err := userEntity.GenerateHashPassword(user.Password())
	if err != nil {
		return userEntity, fmt.Errorf("生成密码失败: %v", err)
	}

	user.UpdatePassword(hash)

	err = us.repo.Save(ctx, user)
	if err != nil {
		return userEntity, fmt.Errorf("保存用户失败: %v", err)
	}
	return userEntity, nil
}

// FindOrCreate 通过手机号查找用户，如果不存在则创建
func (us *userService) FindOrCreate(ctx context.Context, phone string) (*entity.User, error) {
	userEntity := &entity.User{}

	// TODO: 一种优化写法, 大部分情况下都是查找到的
	// 通过手机号查找用户, 如果不存在则创建
	user, err := us.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserDataNotFound) {
		return user, err
	}

	// 注册成功后，再次获取用户信息
	builder := entity.NewUserBuilder()
	userEntity = builder.Phone(phone).Build()

	err = us.repo.Save(ctx, *userEntity)
	if err != nil {
		return userEntity, err
	}

	// TODO: 如果是主从模式下，这里要从主库中读取
	return us.repo.FindByPhone(ctx, phone)
}

func (us *userService) Profile(ctx context.Context, id uint64) (*entity.User, error) {
	return us.repo.FindByID(ctx, id)
}

func (us *userService) UpdateProfile(ctx context.Context, user entity.User) error {
	// 写法1
	// 这种是简单的写法，依赖与 Web 层保证没有敏感数据被修改
	// 也就是说，你的基本假设是前端传过来的数据就是不会修改 Email，Phone 之类的信息的。
	//return svc.repo.Save(ctx, user)

	// 写法2
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 service 层面上维护住了什么是敏感字段这个语义
	user.UpdateEmail("")
	user.UpdatePhone("")
	user.UpdatePassword("")
	err := us.repo.Save(ctx, user)
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}

// UpdateNonSensitiveProfile 更新非敏感数据
// 你可以在这里进一步补充究竟哪些数据会被更新
func (us *userService) UpdateNonSensitiveProfile(ctx context.Context, user entity.User) error {
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 service 层面上维护住了什么是敏感字段这个语义
	user.UpdateEmail("")
	user.UpdatePhone("")
	user.UpdatePassword("")
	return us.repo.Save(ctx, user)
}
