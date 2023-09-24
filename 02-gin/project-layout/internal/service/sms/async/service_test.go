package async

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"project-layout/internal/domain/entity"
	"project-layout/internal/repository"
	"project-layout/internal/repository/mocks"
	"project-layout/internal/service/sms"
	"project-layout/internal/service/sms/mocks"
)

func TestService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, repository.SMSRepository)
		wantErr error
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (sms.Service, repository.SMSRepository) {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "", []string{}, "13888888888")
				return svc, nil
			},
		},
		{
			name: "发送失败",
			mock: func(ctrl *gomock.Controller) (sms.Service, repository.SMSRepository) {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "", []string{}, "13888888888").
					Return(errors.New("发送失败"))
				return svc, nil
			},
			wantErr: errors.New("发送失败"),
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, repository.SMSRepository) {
				svc := smsmocks.NewMockService(ctrl)
				repo := repomocks.NewMockSMSRepository(ctrl)
				svc.EXPECT().Send(gomock.Any(), "", []string{}, "13888888888").
					Return(sms.ErrLimited)
				repo.EXPECT().Save(gomock.Any(), entity.Sms{
					Biz:     "",
					Args:    `[]`,
					Numbers: "13888888888",
					Status:  2,
				})
				return svc, repo
			},
			wantErr: sms.ErrLimited,
		},
		{
			name: "服务商异常",
			mock: func(ctrl *gomock.Controller) (sms.Service, repository.SMSRepository) {
				svc := smsmocks.NewMockService(ctrl)
				repo := repomocks.NewMockSMSRepository(ctrl)
				svc.EXPECT().Send(gomock.Any(), "", []string{}, "13888888888").
					Return(sms.ErrServiceProviderException)
				repo.EXPECT().Save(gomock.Any(), entity.Sms{
					Biz:     "",
					Args:    `[]`,
					Numbers: "13888888888",
					Status:  2,
				})
				return svc, repo
			},
			wantErr: sms.ErrServiceProviderException,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s, repo := tc.mock(ctrl)
			svc := NewService(s, repo, 3)
			err := svc.Send(context.Background(), "", []string{}, "13888888888")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
