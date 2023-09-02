package code

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/redis/go-redis/v9"

	"project-layout/internal/infra"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

type codeRedisCache struct {
	client redis.Cmdable
}

func NewCodeRedisCache(data *infra.Data) CodeCache {
	return &codeRedisCache{
		client: data.RDB,
	}
}

// Set 设置验证码
// 如果该手机再该业务场景下, 验证码不存在(或者已经过期), 则发送验证码
// 如果已经有一个验证码存在, 但是还没有过期, 则返回 ErrUnknownForCode
// 如果已经有一个验证码存在, 且已经过期(验证码已经过去一分钟了), 则重新发送验证码
// 如果已经有一个验证码存在, 但是还没有过期, 不允许重发
// 发送验证码的次数超过了 3 次, 则返回 ErrCodeSendTooMany
// 验证码有效期为 10 分钟
func (r *codeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := r.client.Eval(ctx, luaSetCode, []string{r.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		// 最近一次发送的验证码还没有过期, 不允许重发
		return ErrCodeSendTooMany
	default:
		// 发送验证码遇到未知错误, 系统错误, 比如说 -2, 是key冲突了
		// 其它响应码, 不知道是什么错误

		// TODO:: 这里应该有一个监控, 监控到这里的错误, 然后报警

		return ErrUnknownForCode
	}
}

// Verify 验证验证码是否正确
// 如果验证码是一致的, 那么删除
// 如果验证码不一致, 那么保留
func (r *codeRedisCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := r.client.Eval(ctx, luaVerifyCode, []string{r.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		// 验证次数耗尽, 一般有可能是恶意攻击, 比如说有人在不停的尝试验证码
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不一致
		return false, nil
	}
}

func (r *codeRedisCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
