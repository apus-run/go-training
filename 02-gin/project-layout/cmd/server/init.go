package main

import (
	"os"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"project-layout/internal/service/sms"
	"project-layout/internal/service/sms/localsms"
	"project-layout/internal/service/sms/tencent"
	"project-layout/pkg/ginx"
	"project-layout/pkg/ginx/middleware/auth"
	"project-layout/pkg/ginx/middleware/cors"
	ratelimitRedisMid "project-layout/pkg/ginx/middleware/ratelimit/redis"
	"project-layout/pkg/ratelimit_redis"
)

func InitWebServer(mdls []gin.HandlerFunc, r ginx.Router) *ginx.HttpServer {
	s := ginx.NewHttpServer(
		ginx.WithPort("9000"),
		ginx.WithMode("prod"),
	)
	s.Run(mdls, r)
	return s
}

func InitMiddlewares(client redis.Cmdable) []gin.HandlerFunc {
	rl := ratelimit_redis.NewRedisSlidingWindowLimiter(client, time.Second, 100)
	// 注册中间件
	return []gin.HandlerFunc{
		cors.NewCORS(
			// 允许前端发送
			cors.WithAllowHeaders([]string{"Content-Type", "Authorization"}),
			// 允许前端获取
			cors.WithExposeHeaders([]string{"x-jwt-token"}),
			cors.WithMaxAge(12*60*60),
		).Build(),

		auth.NewBuilder().
			IgnorePaths("user/login").
			IgnorePaths("user/login_sms/code/send").
			IgnorePaths("user/login_sms").
			IgnorePaths("user/register").
			Build(),
		ratelimitRedisMid.NewBuilder(rl).Build(),
	}
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
