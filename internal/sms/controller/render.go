package controller

import (
	"fmt"
	"free/pkg/ext"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func renderError(ctx *gin.Context, err error) {
	cfg := ext.DomainConfig(ctx)

	logrus.WithError(err).WithFields(logrus.Fields{
		"ip":   ctx.ClientIP(),
		"ua":   ctx.Request.UserAgent(),
		"url":  ctx.Request.URL.String(),
		"host": ctx.Request.Host,
	}).Error("server internal error")

	ctx.HTML(200, fmt.Sprintf("%s/%s", cfg.Theme, "error.gohtml"), gin.H{
		"message": err.Error(),
		"router":  ctx.Request.URL.String(),
	})
}

func render(ctx *gin.Context, code int, tpl string, data gin.H) {
	cfg := ext.DomainConfig(ctx)
	data["Host"] = getFullHost(ctx)
	ctx.HTML(code, fmt.Sprintf("%s/%s", cfg.Theme, tpl), data)
}

func getFullHost(ctx *gin.Context) string {
	var (
		host string
	)
	if h := ctx.Request.Host; h != "" {
		host = strings.Split(h, ":")[0]
	}
	if host == "" || host == "127.0.0.1" || host == "localhost" {
		host = ctx.Request.Host
	}
	return fmt.Sprintf("//%s", host)
}
