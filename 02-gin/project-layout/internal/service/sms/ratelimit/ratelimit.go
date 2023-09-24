package ratelimit

import (
	"context"
	"errors"
	"fmt"

	"project-layout/internal/service/sms"
	ratelimit "project-layout/pkg/ratelimit"
)

const SMS_KEY = "sms_tencent"

// Service 定义一个装饰器, 来扩展 SMS Service, 这样就不用改变和污染原有业务逻辑
type Service struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) *Service {
	return &Service{
		svc:     svc,
		limiter: limiter,
	}
}

func (r *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, SMS_KEY)
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流异常 %w", err)
	}
	if limited {
		return errors.New("短信服务触发限流")
	}

	// 以上代码就是用装饰器模式实现限流的短信服务
	// 最终业务逻辑交给了被装饰实现
	return r.svc.Send(ctx, tplId, args, numbers...)
}
