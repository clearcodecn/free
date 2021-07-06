package http

import (
	"context"
	"fmt"
	"free/config"
	"free/internal/sms/controller"
	"free/pkg/ext"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	httpServer *http.Server
)

func Serve() error {
	httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GetConfig().Web.Addr, config.GetConfig().Web.Port),
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logrus.Infof("服务器启动：http://%s:%s", config.GetConfig().Web.Addr, config.GetConfig().Web.Port)
	return httpServer.ListenAndServe()
}

func Stop() error {
	logrus.Info("正在关闭服务器")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	logrus.Info("服务器关闭成功")
	return nil
}

func router() http.Handler {
	r := gin.New()

	r.HTMLRender = ext.NewRender(config.GetConfig().Web.TemplateDir, config.IsProd())

	if err := initRouter(r); err != nil {
		panic(err)
	}

	return r
}

func initRouter(r *gin.Engine) error {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	r.Use(ext.CORS())
	r.Use(ext.Logger())
	r.Use(gin.Recovery())
	r.Use(ext.Cache(ext.CacheConfig{
		Directory: config.GetConfig().Web.CacheDirectory,
		Duration:  config.GetConfig().Web.CacheDuration,
		Enable:    config.GetConfig().Web.EnableCache,
	}))
	r.Use(ext.WithDomainConfig())

	r.GET("/", controller.IndexHandler)
	r.NoRoute(ext.NotFoundHandler())

	ext.GoogleAdsTxt(r)
	ext.Robot(r)
	ext.StaticServer(r, config.GetConfig().Web.StaticDir)
	ext.AddTemporyRouters(r)

	ext.WithDomainConfig()

	return nil
}
