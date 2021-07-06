package controller

import (
	"free/pkg/ext"

	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	d := ext.DomainConfig(ctx)

	render(ctx, 200, "index.gohtml", gin.H{
		"domain": d,
	})
}
