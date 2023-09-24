package alibaba

import (
	"context"
	"os"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	alibabasms "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"

	"project-layout/pkg/utils"
)

func Test_service_Send(t *testing.T) {

	accessKeyId, ok := os.LookupEnv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	require.True(t, ok)
	accessKeySecret, ok := os.LookupEnv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	require.True(t, ok)
	endpoint, ok := os.LookupEnv("ENDPOINT")
	require.True(t, ok)
	signName, ok := os.LookupEnv("SIGN_NAME")
	require.True(t, ok)

	config := &openapi.Config{
		AccessKeyId:     utils.ToPtr[string](accessKeyId),
		AccessKeySecret: utils.ToPtr[string](accessKeySecret),
		Endpoint:        utils.ToPtr[string](endpoint),
	}

	c, err := alibabasms.NewClient(config)
	require.NoError(t, err)
	s := NewService(c, signName)

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:    "发送验证码",
			tplId:   "SMS_154950909",
			params:  []string{"123456"},
			numbers: []string{"15659118048"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}

}
