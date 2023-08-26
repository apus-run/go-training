package tencent

import (
	"context"
	"fmt"
	"project-layout/internal/service/sms"

	tsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"project-layout/pkg/utils"
	"project-layout/pkg/utils/slice"
)

type Service struct {
	client   *tsms.Client
	appId    *string
	signName *string
}

func NewService(c *tsms.Client, appId string,
	signName string) *Service {
	return &Service{
		client:   c,
		appId:    utils.ToPtr[string](appId),
		signName: utils.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := tsms.NewSendSmsRequest()
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.SmsSdkAppId = s.appId
	// ctx 继续往下传
	req.SetContext(ctx)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	req.TemplateId = utils.ToPtr[string](tplId)
	req.SignName = s.signName
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因：%s",
				*status.Code, *status.Message)
		}
	}
	return nil
}
func (s *Service) SendV1(ctx context.Context, tplId string, args []sms.NameArg, numbers ...string) error {
	req := tsms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = utils.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = slice.Map[sms.NameArg, *string](args, func(idx int, src sms.NameArg) *string {
		return &src.Val
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
