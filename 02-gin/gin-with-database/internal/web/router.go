package web

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"gin-with-database/internal/web/handler"
	"gin-with-database/internal/web/middleware/auth"
	ginslog "gin-with-database/internal/web/middleware/slog"
)

func Router(uh *handler.UserHandler) http.Handler {
	log.Printf("load web")

	// Create a slog logger, which:
	//   - Logs to stdout.
	w := os.Stdout
	logger := slog.New(
		slog.NewJSONHandler(
			w,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	logger.WithGroup("http").
		With("environment", "production").
		With("server", "gin/1.9.0").
		With("server_start_time", time.Now()).
		With("gin_mode", gin.EnvGinMode)
	// [SetDefault]还更新了[log]包使用的默认logger
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(
		auth.NewBuilder().
			IgnorePaths("/login").
			IgnorePaths("/signup").
			IgnorePaths("/ping").Build(),
		ginslog.NewBuilder(logger).Build(),
		gin.Recovery(),
	)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "pong",
		})
	})

	engine.POST("/login", uh.Login)
	engine.POST("/signup", uh.Signup)
	engine.GET("/profile", uh.Profile)

	return engine
}
