package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

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
		return userEntity, errors.Wrap(repository.ErrUserDataNotFound, fmt.Sprintf("通过邮箱查找用户失败: %v", err))
	}
	verify := userEntity.VerifyPassword(user.Password, password)
	if !verify {
		return userEntity, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (us *userService) Register(ctx context.Context, user entity.User) (*entity.User, error) {
	userEntity := &entity.User{}
	hash, err := userEntity.GenerateHashPassword(user.Password)
	if err != nil {
		return userEntity, errors.Wrap(err, fmt.Sprintf("生成密码失败: %v", err))
	}
	user.Password = hash
	err = us.repo.Save(ctx, user)
	if err != nil {
		return userEntity, errors.Wrap(err, fmt.Sprintf("保存用户失败: %v", err))
	}
	return &user, nil
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
	userEntity.Phone = phone
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
	u, err := us.repo.FindByID(ctx, user.ID)
	// 中间层 service 尽量不处理 dao 的 error, 直接透传到它的最上层.
	if err != nil && !errors.Is(err, repository.ErrUserDataNotFound) {
		return err
	}

	u.ID = user.ID
	u.Gender = user.Gender
	u.NickName = user.NickName
	u.RealName = user.RealName
	u.Birthday = user.Birthday
	u.Profile = user.Profile

	err = us.repo.Save(ctx, *u)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("更新用户失败: %v", err))
	}
	return nil
}

// UpdateNonSensitiveProfile 更新非敏感数据
// 你可以在这里进一步补充究竟哪些数据会被更新
func (us *userService) UpdateNonSensitiveProfile(ctx context.Context, user entity.User) error {
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 service 层面上维护住了什么是敏感字段这个语义
	user.Email = ""
	user.Phone = ""
	user.Password = ""

	return us.repo.Save(ctx, user)
}
