package service

import (
	"context"

	"github.com/clearcodecn/yunsms/internal/sms/dao"
	"github.com/clearcodecn/yunsms/internal/sms/dao/pg"
	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
)

type MessageService struct {
	dao dao.MessageDao
}

func NewMessageService(domain string) *MessageService {
	return &MessageService{
		dao: pg.NewMessageDao(domain),
	}
}

func (m *MessageService) ListMessageByPhone(ctx context.Context, phone string, pager ext.Pager) ([]*model.Message, int64, error) {
	return m.dao.ListMessageByPhone(ctx, phone, pager)
}
