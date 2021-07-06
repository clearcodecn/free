package ext

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PhoneListPageSize   = 9
	MessageListPageSize = 15
)

type Pager struct {
	// 从1开始
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func Page(ctx *gin.Context, pageSize int) Pager {
	page := ctx.Query("page")
	if page == "" {
		page = "1"
	}
	i, _ := strconv.Atoi(page)
	if i <= 1 {
		i = 1
	}
	return Pager{
		Page:     i,
		PageSize: pageSize,
	}
}

func (p Pager) From() int {
	return (p.Page - 1) * p.PageSize
}

func (p Pager) TotalPage(count int64) int {
	return int(count) / p.PageSize
}
