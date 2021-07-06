package controller

import (
	"free/internal/sms/service"
	"free/pkg/ext"
	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	d := ext.DomainConfig(ctx)

	list := service.NewLinkService().FindAll()

	render(ctx, 200, "index.gohtml", gin.H{
		"domain": d,
		"list":   list,
	})
}
