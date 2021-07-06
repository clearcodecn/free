package ext

import "github.com/gin-gonic/gin"

func Robot(e *gin.Engine) {
	e.GET("/robot.txt", func(ctx *gin.Context) {
		var robot = []byte(`User-agent: YisouSpider
allow: /
User-agent: MipYisouSpider
allow: /
User-agent: Baiduspider
Disallow: /static/
User-agent: Googlebot
Disallow: /static/
User-agent: Sogou web spider
allow: /
User-agent: Sogou inst spider
allow: /
User-agent: Sogou spider2
allow: /
User-agent: Sogou blog
allow: /
User-agent: Sogou News Spider
allow: /
User-agent: Sogou Orion spider
allow: /
User-agent: 360Spider
allow: /
User-agent: Bytespider
Disallow: /
User-agent: *
Disallow: /
`)
		ctx.Writer.Write(robot)
	})
}
