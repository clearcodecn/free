package ext

import "github.com/gin-gonic/gin"

func GoogleAdsTxt(e *gin.Engine) {
	e.GET("/ads.txt", func(ctx *gin.Context) {
		d := DomainConfig(ctx)
		ctx.Writer.Header().Set("Content-Type", "text/plain")
		ctx.Writer.Write([]byte(d.AdsTxt))
	})
}
