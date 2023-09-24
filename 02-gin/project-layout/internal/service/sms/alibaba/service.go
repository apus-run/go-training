package alibaba

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/alibabacloud-go/dysmsapi-20170525/v3/client"

	"project-layout/internal/service/sms"
	"project-layout/pkg/utils"
)

type service struct {
	client   *client.Client
	signName *string
}

func NewService(client *client.Client, signName string) sms.Service {
	return &service{
		client:   client,
		signName: utils.ToPtr[string](signName),
	}
}

func (s *service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := &client.SendSmsRequest{
		SignName:      s.signName,
		TemplateCode:  utils.ToPtr[string]("SMS_154950909"),
		PhoneNumbers:  utils.ToPtr[string](strings.Join(numbers, ",")),
		TemplateParam: utils.ToPtr[string](fmt.Sprintf(`{"code":"%s"}`, args[0])),
	}

	resp, err := s.client.SendSms(req)
	if err != nil {
		log.Println("发送短信失败:", err)
		return err
	}
	if *(resp.Body.Code) != "OK" {
		log.Println(resp.Body.String())
	}
	return nil
}
