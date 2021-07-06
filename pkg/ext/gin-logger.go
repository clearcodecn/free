package ext

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 生产环境下，打印req body 但是不打印 res body
	return func(c *gin.Context) {
		reqid := setRequestIDHeader(c)

		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		// 执行时间
		rt := end.Sub(start)

		// 请求方式
		method := c.Request.Method

		// 请求路由
		url := c.Request.URL.String()

		// 状态码
		code := c.Writer.Status()

		// 请求IP
		client := c.ClientIP()

		// 协议
		proto := c.Request.Proto

		logger.WithFields(logrus.Fields{
			"type":            "http",
			"code":            code,
			"rt":              rt,
			"client":          client,
			"method":          method,
			"url":             url,
			"proto":           proto,
			"reqid":           reqid,
			"body_bytes_sent": c.Request.ContentLength,
			"host":            c.Request.Host,
		}).Info()
	}
}

const (
	HeaderXRequestId = "X-Request-Id"
)

// setRequestIDHeader 填充 requestid 到header中
func setRequestIDHeader(c *gin.Context) string {
	if c.GetHeader(HeaderXRequestId) != "" {
		return c.GetHeader(HeaderXRequestId)
	}
	id := uuid.New().String()
	c.Request.Header.Set(HeaderXRequestId, id)
	return id
}
