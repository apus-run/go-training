package sms

import (
	"context"
)

// Service 发送短信的抽象
// 目前你可以理解为，这是一个为了适配不同的短信供应商的抽象
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error

	SendV1(ctx context.Context, tplId string, args []NameArg, numbers ...string) error
}

// NameArg 短信模板中的参数, 例如: 您的验证码是: {{code}} , 那么这里的 Name 就是 code, Val 就是验证码的值
type NameArg struct {
	Val  string
	Name string
}
