package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"project-layout/internal/domain/entity"
	"project-layout/internal/infra"
)

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotExist = redis.Nil

var _ UserCache = (*userCache)(nil)

type UserCache interface {
	// Set 理论上来说，UserCache 也应该有自己的 User 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, user entity.User) error
	Get(ctx context.Context, id uint64) (entity.User, error)

	Remove(ctx context.Context, id uint64) error
}

type userCache struct {
	data *infra.Data

	// 缓存过期时间
	expire time.Duration
}

func NewUserCache(data *infra.Data) UserCache {
	return &userCache{
		data:   data,
		expire: time.Minute * 15,
	}
}

func (u *userCache) Set(ctx context.Context, user entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := u.key(user.ID)
	return u.data.RDB.Set(ctx, key, string(data), u.expire).Err()
}

func (u *userCache) Get(ctx context.Context, id uint64) (entity.User, error) {
	key := u.key(id)
	data, err := u.data.RDB.Get(ctx, key).Bytes()
	if err != nil {
		return entity.User{}, err
	}

	// 反序列化回来
	var art entity.User
	err = json.Unmarshal(data, &art)
	return art, err
}

func (u *userCache) Remove(ctx context.Context, id uint64) error {
	key := u.key(id)
	return u.data.RDB.Del(ctx, key).Err()
}

func (u *userCache) key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}
