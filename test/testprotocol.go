package main

import "github.com/gin-gonic/gin"

func main() {
	g := gin.Default()

	g.GET("/", func(context *gin.Context) {
		header := context.Request.Header

		h := gin.H{}
		h["header"] = header

		h["proto"] = context.Request.Proto
		h["ProtoMajor"] = context.Request.ProtoMajor
		h["ProtoMinor"] = context.Request.ProtoMinor
		h["host"] = context.Request.Host
		h["remoteAddr"] = context.Request.RemoteAddr
		h["requestURI"] = context.Request.RequestURI
		url := context.Request.URL
		h["url"] = url

		context.JSON(200, h)
	})

	g.Run(":9123")
}
