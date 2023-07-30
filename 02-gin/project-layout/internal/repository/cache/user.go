package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
)

var _ UserCache = (*userCache)(nil)
var ErrNotFound = errors.New("插入失败")

type UserCache interface {
	// Set 理论上来说，UserCache 也应该有自己的 User 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, user entity.User) error
	Get(ctx context.Context, id uint64) (entity.User, error)
}

type userCache struct {
	data *repository.Data
}

func (u *userCache) Set(ctx context.Context, user entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	res, err := u.data.RDB.Set(ctx, fmt.Sprintf("article_%d", user.ID), string(data), time.Hour).Result()
	if res != "OK" {
		return ErrNotFound
	}
	return err
}

func (u *userCache) Get(ctx context.Context, id uint64) (entity.User, error) {
	data, err := u.data.RDB.Get(ctx, fmt.Sprintf("article_%d", id)).Bytes()
	if err != nil {
		return entity.User{}, err
	}
	var art entity.User
	err = json.Unmarshal(data, &art)
	return art, err
}
