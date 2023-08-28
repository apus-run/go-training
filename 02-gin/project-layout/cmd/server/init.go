package main

import (
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"project-layout/internal/service/sms"
	"project-layout/internal/service/sms/localsms"
	"project-layout/internal/service/sms/tencent"
	"project-layout/pkg/ginx"
)

func InitWebServer(r ginx.Router) *ginx.HttpServer {
	s := ginx.NewHttpServer(
		ginx.WithPort("9000"),
		ginx.WithMode("prod"),
	)
	s.Run(r)
	return s
}

func InitSmsService() sms.Service {
	//return initSmsTencentService()
	return initSmsMemoryService()
}

func initSmsTencentService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("没有找到环境变量 SMS_SECRET_ID ")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	c, err := tencentSMS.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic("没有找到环境变量 SMS_SECRET_KEY")
	}
	return tencent.NewService(c, "15800000008", "无码科技")
}

// initSmsMemoryService 方便测试, 使用基于内存，输出到控制台的实现, 模拟短信服务
func initSmsMemoryService() sms.Service {
	return localsms.NewService()
}
