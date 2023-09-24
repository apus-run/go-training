package failover

import (
	"context"
	"log"
	"sync/atomic"

	"project-layout/internal/service/sms"
)

type service struct {
	svcs []sms.Service
	idx  uint64
}

func NewService(svcs []sms.Service) sms.Service {
	return &service{svcs: svcs}
}

func (s *service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&s.idx, 1)
	length := uint64(len(s.svcs))
	for i := idx; i < idx+length; i++ {
		svc := s.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 调用者设置的超时时间到了
			// 调用者主动取消了
			return err
		default:
			log.Printf("发送失败: %v", svc)
		}
	}
	return sms.ErrServiceProviderException
}

func (s *service) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range s.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Printf("发送失败: %v", svc)
	}
	return sms.ErrServiceProviderException
}
