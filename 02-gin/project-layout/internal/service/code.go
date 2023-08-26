package service

import (
	"context"
	"fmt"
	"math/rand"

	"project-layout/internal/repository"
	"project-layout/internal/service/sms"
)

var (
	codeTplId                 = "1933434"
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// codeService 短信验证码服务
type codeService struct {
	sms  sms.Service
	repo repository.CodeRepository
}

// NewCodeService 创建一个验证码服务
func NewCodeService(sms sms.Service, repo repository.CodeRepository) CodeService {
	return &codeService{
		sms:  sms,
		repo: repo,
	}
}

// Send 发送验证码, biz 区别业务场景, mobile 是手机号, code 是验证码
func (s *codeService) Send(ctx context.Context, biz, phone string) error {
	// 生成一个验证码
	code := s.generateCode()

	// 放到Redis中
	err := s.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	// 发送验证码
	err = s.sms.Send(ctx, codeTplId, []string{code}, phone)

	// TODO:: 如果err不为nil, 要处理err错误吗? 例如: 重试, 或者记录日志
	// TODO:: 如果err不为nil, 要把Redis中的验证码删除吗?

	// 原样返回err
	return err
}

func (s *codeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, code)

	// TODO:: 我们可以在这里屏蔽一些错误, 如果验证码验证次数过多(有可能是恶意攻击), 我们可以屏蔽这个错误, 就不透传到上层handler了
	if err == ErrCodeVerifyTooManyTimes {
		// TODO:: 这里可以记录日志, 例如: 有人在恶意攻击我们的短信验证码

		return false, err
	}
	return ok, nil
}

func (s *codeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%6d", num)
}
