package ginx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"project-layout/pkg/log"
)

// HttpServer 代表当前服务端实例
type HttpServer struct {
	// 服务配置
	options *Options
	// 日志
	log *log.Logger
	// shutdown回调函数
	fn func()
}

// NewHttpServer 创建server实例
func NewHttpServer(logger *log.Logger, options ...Option) *HttpServer {
	opts := Apply(options...)
	return &HttpServer{
		options: opts,
		log:     logger,
	}
}

// Router 加载路由，使用侧提供接口，实现侧需要实现该接口
type Router interface {
	Load(engine *gin.Engine)
}

// Run server的启动入口
// 加载路由, 启动服务
func (s *HttpServer) Run(rs ...Router) {
	var wg sync.WaitGroup
	wg.Add(1)

	// 设置gin启动模式，必须在创建gin实例之前
	gin.SetMode(s.options.mode)

	// 创建gin实例
	g := gin.Default()

	// 加载路由
	s.registerRoutes(g, rs...)

	// health check
	go func() {
		if err := Ping(s.options.port, s.options.maxPingCount); err != nil {
			s.log.Fatal("server no response")
		}
		s.log.Infof("server started success! port: %s", s.options.port)
	}()

	srv := http.Server{
		Addr:    s.options.port,
		Handler: g,
	}
	if s.fn != nil {
		srv.RegisterOnShutdown(s.fn)
	}
	// graceful shutdown
	sgn := make(chan os.Signal, 1)
	signal.Notify(
		sgn,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	go func() {
		<-sgn
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			s.log.Errorf("server shutdown err %v \n", err)
		}
		wg.Done()
	}()

	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			s.log.Errorf("server start failed on port %s", s.options.port)
			return
		}
	}
	wg.Wait()
	s.log.Infof("server stop on port %s", s.options.port)
}

// RouterLoad 加载自定义路由
func (s *HttpServer) registerRoutes(g *gin.Engine, rs ...Router) *HttpServer {
	for _, r := range rs {
		r.Load(g)
	}
	return s
}

// RegisterOnShutdown 注册shutdown后的回调处理函数，用于清理资源
func (s *HttpServer) RegisterOnShutdown(f func()) {
	s.fn = f
}

// Ping 用来检查是否程序正常启动
func Ping(port string, maxCount int) error {
	seconds := 1
	if len(port) == 0 {
		panic("Please specify the service port")
	}
	if !strings.HasPrefix(port, ":") {
		port += ":"
	}
	url := fmt.Sprintf("http://localhost%s/ping", port)
	for i := 0; i < maxCount; i++ {
		resp, err := http.Get(url)
		if nil == err && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		fmt.Printf("等待服务在线, 已等待 %d 秒，最多等待 %d 秒", seconds, maxCount)
		time.Sleep(time.Second * 1)
		seconds++
	}
	return fmt.Errorf("服务启动失败，端口 %s", port)
}
