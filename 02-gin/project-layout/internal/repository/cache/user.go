package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"project-layout/internal/domain/entity"
	"project-layout/internal/infra"
)

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotFound = errors.New("key not found")

var _ UserCache = (*userCache)(nil)

type UserCache interface {
	// Set 理论上来说，UserCache 也应该有自己的 User 定义
	// 比如说你并不需要缓存全部字段
	// 但是我们这里直接缓存全部
	Set(ctx context.Context, key string, val string) error
	SetOjb(ctx context.Context, user entity.User) error
	Get(ctx context.Context, id uint64) (entity.User, error)
	GetObj(ctx context.Context, key string, model interface{}) error

	Del(ctx context.Context, id uint64) error
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

func (uc *userCache) Set(ctx context.Context, key string, val string) error {
	return uc.data.RDB.Set(ctx, key, val, uc.expire).Err()
}

// SetObj 设置某个key和对象到缓存, 对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (uc *userCache) SetOjb(ctx context.Context, user entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := uc.key(user.ID)
	return uc.data.RDB.Set(ctx, key, string(data), uc.expire).Err()
}

// SetMany 设置多个key和值到缓存
func (uc *userCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := uc.data.RDB.Pipeline()
	cmds := make([]*redis.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (uc *userCache) Get(ctx context.Context, id uint64) (entity.User, error) {
	key := uc.key(id)
	cmd := uc.data.RDB.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return entity.User{}, ErrKeyNotFound
	}

	data, _ := cmd.Bytes()

	// 反序列化回来
	var art entity.User
	err := json.Unmarshal(data, &art)
	return art, err
}

// GetObj 获取某个key对应的对象, 对象必须实现 https://pkg.go.dev/encoding#BinaryUnMarshaler
func (uc *userCache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := uc.data.RDB.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return ErrKeyNotFound
	}

	err := cmd.Scan(model)
	if err != nil {
		return err
	}
	return nil
}

// GetMany 获取某些key对应的值
func (uc *userCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := uc.data.RDB.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redis.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	return vals, nil
}

func (uc *userCache) Del(ctx context.Context, id uint64) error {
	key := uc.key(id)
	return uc.data.RDB.Del(ctx, key).Err()
}

func (uc *userCache) DelMany(ctx context.Context, keys []string) error {
	pipline := uc.data.RDB.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (uc *userCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return uc.data.RDB.IncrBy(ctx, key, step).Result()
}

func (uc *userCache) Increment(ctx context.Context, key string) (int64, error) {
	return uc.data.RDB.IncrBy(ctx, key, 1).Result()
}

func (uc *userCache) Decrement(ctx context.Context, key string) (int64, error) {
	return uc.data.RDB.IncrBy(ctx, key, -1).Result()
}

// SetTTL 设置某个key的超时时间
func (uc *userCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return uc.data.RDB.Expire(ctx, key, timeout).Err()
}

// GetTTL 获取某个key的超时时间
func (uc *userCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return uc.data.RDB.TTL(ctx, key).Result()
}

func (u *userCache) key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}
