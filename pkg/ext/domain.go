package ext

import (
	"errors"
	"free/config"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func WithDomainConfig() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		host := ctx.Request.Host
		if strings.Contains(host, ":") {
			host, _, _ = net.SplitHostPort(host)
		}
		if domain := config.GetConfig().GetConfigByHost(host); domain == nil {
			logrus.Errorf("bad host: %v", host)
			ctx.AbortWithError(404, errors.New("bad request ^~^"))
			return
		} else {
			ctx.Set("domain", domain)
			ctx.Next()
		}
	}
}

func DomainConfig(ctx *gin.Context) *config.Domain {
	v, ok := ctx.Get("domain")
	if ok {
		return v.(*config.Domain)
	}
	return nil
}
