package ext

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// 临时的路由.
var temporyRouters = make(map[string]string)
var mux sync.RWMutex

func NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mux.RLock()
		if v, ok := temporyRouters[ctx.Request.RequestURI]; ok {
			mux.RUnlock()

			ctx.Header("Content-Type", "text/plain")
			ctx.String(200, v)
			return
		}

		mux.RUnlock()
		ctx.AbortWithStatus(404)
	}
}

func AddTemporyRouters(r *gin.Engine) {
	fn := func(ctx *gin.Context) {
		key, ok := ctx.GetPostForm("key")
		if !ok {
			ctx.String(200, "key not available")
			return
		}
		val, ok := ctx.GetPostForm("value")
		if !ok {
			ctx.String(200, "value not available")
			return
		}

		mux.Lock()
		defer mux.Unlock()

		temporyRouters[key] = val
	}

	r.Any("/temp", fn)
}
