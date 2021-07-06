package controller

import (
	"free/internal/sms/service"
	"free/pkg/ext"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func DetailHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	d := ext.DomainConfig(ctx)

	list, err := service.NewLinkService().FindOne(id)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read list")
		return
	}

	render(ctx, 200, "index.gohtml", gin.H{
		"domain": d,
		"list":   list,
	})
}
