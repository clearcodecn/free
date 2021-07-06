package dao

import (
	"context"

	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
)

type MessageDao interface {
	ListMessageByPhone(ctx context.Context, phone string, pager ext.Pager) ([]*model.Message, int64, error)
	CreateMessage(ctx context.Context, sms *model.Message) error
	BulkCreateMessage(ctx context.Context, sms []*model.Message) error
	CountByPhone(ctx context.Context, phone string) (int64, error)
}

type PhoneDao interface {
	ListPhone(ctx context.Context, pager ext.Pager, country string) ([]*model.Phone, int64, error)
	GetPhoneById(ctx context.Context, id string) (*model.Phone, error)
	CreatePhone(ctx context.Context, sms *model.Phone) error
	BulkCreatePhone(ctx context.Context, sms []*model.Phone) error
	RandomGet(ctx context.Context, size int) ([]*model.Phone, error)
	UpdateLastOne(ctx context.Context, phone string) error
}
