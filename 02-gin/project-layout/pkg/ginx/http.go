package ginx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// HttpServer 代表当前服务端实例
type HttpServer struct {
	// 服务配置
	options *Options

	// shutdown回调函数
	fn func()
}

// NewHttpServer 创建server实例
func NewHttpServer(options ...Option) *HttpServer {
	opts := Apply(options...)
	return &HttpServer{
		options: opts,
	}
}

// Router 加载路由，使用侧提供接口，实现侧需要实现该接口
type Router interface {
	Load(engine *gin.Engine)
}

// Run server的启动入口
// 加载路由, 启动服务
func (s *HttpServer) Run(mdls []gin.HandlerFunc, rs ...Router) {
	wg := sync.WaitGroup{}

	// 设置gin启动模式，必须在创建gin实例之前
	switch s.options.mode {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 创建gin实例
	g := gin.Default()

	// 注册中间件
	g.Use(mdls...)

	// 加载路由
	s.registerRoutes(g, rs...)

	var addr string
	if strings.HasPrefix(s.options.port, ":") {
		addr = fmt.Sprintf("%s%s", s.options.host, s.options.port)
	}
	addr = fmt.Sprintf("%s:%s", s.options.host, s.options.port)

	// 启动服务
	srv := http.Server{
		Addr:    addr,
		Handler: g,
	}
	if s.fn != nil {
		srv.RegisterOnShutdown(s.fn)
	}
	// graceful shutdown
	exitSignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, len(exitSignals))
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, exitSignals...)

	wg.Add(1)
	go func() {
		<-quit
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(stopCtx); err != nil {
			log.Printf("server shutdown err %v \n", err)
		}
		wg.Done()
	}()

	log.Printf("server start on port %s", s.options.port)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server start failed on port %s", s.options.port)
	}
	wg.Wait()
	log.Printf("server stop on port %s", s.options.port)
}

// RouterLoad 加载自定义路由
func (s *HttpServer) registerRoutes(g *gin.Engine, rs ...Router) *HttpServer {
	for _, r := range rs {
		r.Load(g)
	}
	return s
}

// RegisterOnShutdown 注册shutdown后的回调处理函数，用于清理资源
func (s *HttpServer) RegisterOnShutdown(fn func()) {
	s.fn = fn
}
