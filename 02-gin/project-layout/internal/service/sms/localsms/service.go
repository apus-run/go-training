package localsms

import (
	"context"
	"log"

	"project-layout/internal/service/sms"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是:", args)
	return nil
}

func (s *Service) SendV1(ctx context.Context, tplId string, args []sms.NameArg, numbers ...string) error {
	log.Println("验证码是:", args)
	return nil
}
