package failover

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"project-layout/internal/service/sms"
	"project-layout/internal/service/sms/mocks"
)

func TestService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr error
	}{
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
		},
		{
			name: "重试成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送不了"))
				return []sms.Service{svc0, svc1}
			},
		},
		{
			name: "最终失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("发送不了"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("还是失败"))
				return []sms.Service{svc0, svc1}
			},
			wantErr: errors.New("发送失败，所有服务商都尝试过了"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "mytpl", []string{"123456"}, "13888888888")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
