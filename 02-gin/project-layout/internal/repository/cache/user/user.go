package cache

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
	// 使用 Cmdable 是为了更好的扩展性, redis.Client 和 redis.ClusterClient 都实现了这个接口
	client redis.Cmdable

	// 缓存过期时间
	expire time.Duration
}

func NewUserCache(data *infra.Data) UserCache {
	return &userCache{
		client: data.RDB,
		expire: time.Minute * 15,
	}
}

func (cache *userCache) Set(ctx context.Context, key string, val string) error {
	return cache.client.Set(ctx, key, val, cache.expire).Err()
}

// SetObj 设置某个key和对象到缓存, 对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (cache *userCache) SetOjb(ctx context.Context, user entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.ID())
	return cache.client.Set(ctx, key, string(data), cache.expire).Err()
}

// SetMany 设置多个key和值到缓存
func (cache *userCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := cache.client.Pipeline()
	cmds := make([]*redis.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (cache *userCache) Get(ctx context.Context, id uint64) (entity.User, error) {
	key := cache.key(id)
	cmd := cache.client.Get(ctx, key)
	// 数据不存在，cmd.Err = redis.Nil
	if errors.Is(cmd.Err(), ErrKeyNotExist) {
		return entity.User{}, ErrKeyNotFound
	}

	// 反序列化回来
	var art entity.User
	data, _ := cmd.Bytes()
	err := json.Unmarshal(data, &art)
	return art, err
}

// GetObj 获取某个key对应的对象, 对象必须实现 https://pkg.go.dev/encoding#BinaryUnMarshaler
func (cache *userCache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := cache.client.Get(ctx, key)
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
func (cache *userCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := cache.client.Pipeline()
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

func (cache *userCache) Del(ctx context.Context, id uint64) error {
	key := cache.key(id)
	return cache.client.Del(ctx, key).Err()
}

func (cache *userCache) DelMany(ctx context.Context, keys []string) error {
	pipline := cache.client.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (cache *userCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return cache.client.IncrBy(ctx, key, step).Result()
}

func (cache *userCache) Increment(ctx context.Context, key string) (int64, error) {
	return cache.client.IncrBy(ctx, key, 1).Result()
}

func (cache *userCache) Decrement(ctx context.Context, key string) (int64, error) {
	return cache.client.IncrBy(ctx, key, -1).Result()
}

// SetTTL 设置某个key的超时时间
func (cache *userCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return cache.client.Expire(ctx, key, timeout).Err()
}

// GetTTL 获取某个key的超时时间
func (cache *userCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return cache.client.TTL(ctx, key).Result()
}

func (cache *userCache) key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}
