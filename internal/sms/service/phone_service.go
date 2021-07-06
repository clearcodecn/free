package service

import (
	"context"

	"github.com/clearcodecn/yunsms/internal/sms/dao"
	"github.com/clearcodecn/yunsms/internal/sms/dao/pg"
	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
)

type PhoneService struct {
	dao dao.PhoneDao
}

func (s *PhoneService) ListPhone(ctx context.Context, pager ext.Pager, country string) ([]*model.Phone, int64, error) {
	return s.dao.ListPhone(ctx, pager, country)
}

func (s *PhoneService) GetPhoneByID(ctx context.Context, id string) (*model.Phone, error) {
	return s.dao.GetPhoneById(ctx, id)
}

func NewPhoneService(db string) *PhoneService {
	return &PhoneService{
		dao: pg.NewPhoneDao(db),
	}
}
