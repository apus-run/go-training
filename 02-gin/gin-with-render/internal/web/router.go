package web

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hako/durafmt"

	"gin-with-render/internal/web/handler"
	"gin-with-render/internal/web/middleware/auth"
)

var (
	//go:embed assets/*
	staticFiles embed.FS

	//go:embed templates/*
	templateFiles embed.FS
)

func Router(uh *handler.UserHandler) http.Handler {
	log.Printf("load web")

	// gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(
		auth.NewBuilder().
			IgnorePaths("/favicon.ico").
			IgnorePaths("/public").
			IgnorePaths("/foo.html").
			IgnorePaths("/index.html").
			IgnorePaths("/login").
			IgnorePaths("/signup").
			IgnorePaths("/ping").Build(),
		gin.Recovery(),
	)

	templ := template.Must(
		template.New("").
			Funcs(template.FuncMap{
				"niceSizeMB": func(s int) string { return fmt.Sprintf("%.1f", float32(s)/1024/1024) },
				"join":       strings.Join,
			}).
			ParseFS(
				templateFiles,
				"templates/*.html",
				"templates/foo/*.html",
			))
	engine.SetHTMLTemplate(templ)

	//images, err := fs.Sub(staticFiles, "images")
	//if err != nil {
	//	log.Fatalf("problem with assetFS: %s", err)
	//}
	// example: /public/assets/images/example.png
	engine.StaticFS("/public", http.FS(staticFiles))

	engine.GET("favicon.ico", func(c *gin.Context) {
		file, _ := staticFiles.ReadFile("assets/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})

	engine.GET("/foo.html", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"bar.html",
			gin.H{
				"title": "Foo website",
				"name":  "About",
				"msg":   "All about Charlie",
			})
	})

	engine.GET("/index.html", func(ctx *gin.Context) {
		ctx.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "hello html",
				"name":  "HOME",
				"msg":   "Hello, Gent!",
			},
		)
	})

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
