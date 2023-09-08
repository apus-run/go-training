package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"project-layout/internal/domain/entity"
	"project-layout/internal/infra"
)

var ErrKeyNotFound = errors.New("key not found")

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotExist = redis.Nil

var _ UserCache = (*userRedisCache)(nil)

type UserCache interface {
	// Set 理论上来说，UserCache 也应该有自己的 User 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, key string, val string) error
	SetOjb(ctx context.Context, user entity.User) error
	Get(ctx context.Context, key string) (string, error)
	GetObj(ctx context.Context, id uint64) (entity.User, error)

	Del(ctx context.Context, id uint64) error
	DelMany(ctx context.Context, keys []string) error
}

type userRedisCache struct {
	// 使用 Cmdable 是为了更好的扩展性, redis.Client 和 redis.ClusterClient 都实现了这个接口
	client redis.Cmdable

	// 缓存过期时间
	expire time.Duration
}

func NewUserRedisCache(data *infra.Data) UserCache {
	return &userRedisCache{
		client: data.RDB,
		expire: time.Minute * 15,
	}
}

func (cache *userRedisCache) Set(ctx context.Context, key string, val string) error {
	return cache.client.Set(ctx, key, val, cache.expire).Err()
}

// SetObj 设置某个key和对象到缓存中
func (cache *userRedisCache) SetOjb(ctx context.Context, user entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.ID)
	return cache.client.Set(ctx, key, string(data), cache.expire).Err()
}

func (cache *userRedisCache) Get(ctx context.Context, key string) (string, error) {
	cmd := cache.client.Get(ctx, key)
	// 数据不存在，cmd.Err = redis.Nil
	if errors.Is(cmd.Err(), ErrKeyNotExist) {
		return "", ErrKeyNotFound
	}
	return cmd.Result()
}

// GetObj 获取某个key对应的对象
func (cache *userRedisCache) GetObj(ctx context.Context, id uint64) (entity.User, error) {
	key := cache.key(id)

	cmd := cache.client.Get(ctx, key)
	if errors.Is(cmd.Err(), ErrKeyNotExist) {
		return entity.User{}, ErrKeyNotFound
	}

	// 反序列化回来
	var u entity.User
	data, _ := cmd.Bytes()
	err := json.Unmarshal(data, &u)
	return u, err
}

func (cache *userRedisCache) Del(ctx context.Context, id uint64) error {
	key := cache.key(id)
	return cache.client.Del(ctx, key).Err()
}

func (cache *userRedisCache) DelMany(ctx context.Context, keys []string) error {
	pipline := cache.client.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (cache *userRedisCache) key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}
