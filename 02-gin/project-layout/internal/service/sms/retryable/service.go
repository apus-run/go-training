package retryable

import (
	"context"
	"errors"
	"sync"

	"project-layout/internal/service/sms"
)

type Service struct {
	svc sms.Service

	// 重试
	retryMax int

	mu sync.Mutex
}

func NewService(svc sms.Service) *Service {
	return &Service{
		svc:      svc,
		retryMax: 3,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.Send(ctx, tpl, args, numbers...)
	cnt := 1

	for err != nil && cnt < s.retryMax {
		err = s.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}

	return errors.New("重新都失败了")
}
