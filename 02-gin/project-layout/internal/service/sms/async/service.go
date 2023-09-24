package async

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
	"project-layout/internal/service/sms"
	"strings"
	"sync"
	"time"
)

type Service struct {
	svc  sms.Service
	repo repository.SMSRepository

	retryMax int // 重试次数
}

func NewService(svc sms.Service, repo repository.SMSRepository, retryMax int) *Service {
	return &Service{
		svc:  svc,
		repo: repo,

		retryMax: retryMax,
	}
}

func (s *Service) Async(ctx context.Context) {
	var wg sync.WaitGroup

	requests, err := s.repo.FindByStatus(ctx, 2)
	if err != nil {
		fmt.Printf("record not found")
	}
	for _, request := range requests {
		wg.Add(1)
		go func(sms entity.Sms) {
			defer wg.Done()

			// 发送请求并进行重试
			for i := 0; i < s.retryMax; i++ {
				args := strings.Split(sms.Args, ";")
				numbers := strings.Split(sms.Numbers, ",")
				err := s.svc.Send(ctx, sms.Biz, args, numbers...)
				if err != nil {
					fmt.Printf("Failed to send request: %v\n", err)
				}
				if err == nil {
					// 请求发送成功, 退出重试循环
					return
				}

				// 重试间隔为 1 秒
				time.Sleep(time.Second)
			}

		}(request)
	}

	// 等待所有请求处理完成
	wg.Wait()

	// 休眠一段时间后再次检查数据库中是否有新的请求
	time.Sleep(time.Second)
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 服务正常
	err := s.svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		// 判定是否限流 或者 崩溃了
		if errors.Is(err, sms.ErrLimited) || errors.Is(err, sms.ErrServiceProviderException) {
			bs, er := json.Marshal(&args)
			if er != nil {
				return er
			}
			s.repo.Save(ctx, entity.Sms{
				Biz:     tplId,
				Args:    string(bs),
				Numbers: strings.Join(numbers, ","),
				Status:  2,
			})
		}
	}
	return nil
}
