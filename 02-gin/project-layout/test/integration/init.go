package integration

import (
	"os"

	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"project-layout/internal/infra"
	"project-layout/internal/repository"
	"project-layout/internal/repository/cache/code"
	"project-layout/internal/repository/cache/user"
	"project-layout/internal/repository/dao"
	"project-layout/internal/service"
	"project-layout/internal/service/oauth2/wechat"
	"project-layout/internal/service/sms"
	"project-layout/internal/service/sms/localsms"
	"project-layout/internal/service/sms/tencent"
	"project-layout/internal/web/handler"
	"project-layout/internal/web/handler/jwt"
	"project-layout/pkg/conf"
	"project-layout/pkg/conf/file"
	"project-layout/pkg/ginx/middleware/auth"
	"project-layout/pkg/ginx/middleware/cors"
	ratelimitRedisMid "project-layout/pkg/ginx/middleware/ratelimit/redis"
	"project-layout/pkg/log"
	"project-layout/pkg/ratelimit"
)

func initWebServer(mdls []gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Use(mdls...)
	return r
}

func initMiddlewares(client redis.Cmdable) []gin.HandlerFunc {
	rl := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 100)
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

func initSmsService() sms.Service {
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

func runApp(logger *log.Logger) *gin.Engine {
	cmdable := infra.NewRDB()
	db := infra.NewDB()
	data, _ := infra.NewData(db, cmdable, logger)
	userDAO := dao.NewUserDAO(data)
	userCache := user.NewUserRedisCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache, logger)
	userService := service.NewUserService(userRepository, logger)
	smsService := initSmsService()
	codeCache := code.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	codeService := service.NewCodeService(smsService, codeRepository)
	jwtHandler := jwt.NewJwtHandler()
	userHandler := handler.NewUserHandler(userService, codeService, jwtHandler, logger)
	wechatService := wechat.NewService()
	oauth2WechatHandler := handler.NewOAuth2WechatHandler(wechatService, userService, jwtHandler)
	v := initMiddlewares(cmdable)
	server := initWebServer(v)
	userHandler.Load(server)
	oauth2WechatHandler.Load(server)

	return server
}

func Start() *gin.Engine {
	c := conf.New(
		conf.WithSource(
			file.NewSource("./config"),
		),
	)
	c.Load()
	c.Watch()

	logger := log.NewLogger(
		log.WithEncoding("json"),
		log.WithFilename("./logs/server.log"),
	)

	defer logger.Sync()

	s := runApp(logger)

	return s
}
