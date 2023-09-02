package code

import (
	"context"
	"errors"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码过于频繁")
	ErrUnknownForCode         = errors.New("发送验证码遇到未知错误")
	ErrCodeVerifyTooManyTimes = errors.New("验证码验证次数过多")
)

type CodeCache interface {
	// Set 设置验证码, biz 为业务类型, 比如注册, 登录等
	Set(ctx context.Context, biz, phone, code string) error

	// Verify 验证验证码是否正确, biz 为业务类型, 比如注册, 登录等
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}
