package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/clearcodecn/yunsms/internal/sms/dao"
	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

type phoneDao struct {
	table string
}

func (p *phoneDao) ListPhone(ctx context.Context, pager ext.Pager, country string) ([]*model.Phone, int64, error) {
	sess := p.newSession(phoneTable)
	defer sess.Close()

	var list []*model.Phone
	if country != "" {
		sess.Where("flag = ?", country)
	}
	sess.OrderBy("last_one desc")
	sess.Limit(pager.PageSize, pager.From())

	count, err := sess.FindAndCount(&list)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (p *phoneDao) GetPhoneById(ctx context.Context, id string) (*model.Phone, error) {
	sess := p.newSession(phoneTable)
	defer sess.Close()

	var phone = new(model.Phone)
	ok, err := sess.Where("id = ?", id).Get(phone)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sql.ErrNoRows
	}
	return phone, nil
}

func (p *phoneDao) CreatePhone(ctx context.Context, sms *model.Phone) error {
	sess := p.newSession(phoneTable)
	defer sess.Close()
	_, err := sess.Insert(sms)
	return err
}

func (p *phoneDao) BulkCreatePhone(ctx context.Context, sms []*model.Phone) error {
	sess := p.newSession(phoneTable)
	defer sess.Close()

	for {
		var data []*model.Phone
		if len(sms) > 10000 {
			data = sms[:10000]
			sms = sms[10000:]
		} else {
			data = sms
			sms = sms[:0]
		}
		if _, err := p.setTable(sess, phoneTable).InsertMulti(data); err != nil {
			return err
		}
		if len(sms) == 0 {
			break
		}
	}

	return nil
}

func (p *phoneDao) RandomGet(ctx context.Context, size int) ([]*model.Phone, error) {
	sess := p.newSession(phoneTable)
	defer sess.Close()

	var list []*model.Phone
	err := sess.SQL(fmt.Sprintf("select * from %s order by random() limit ?", p.table+"_"+phoneTable), size).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (p *phoneDao) UpdateLastOne(ctx context.Context, phone string) error {
	sess := p.newSession(phoneTable)
	defer sess.Close()

	_, err := sess.Where("phone_no = ?", phone).Cols("last_one").Update(&model.Phone{
		LastOne: time.Now(),
	})
	return err
}

func (p *phoneDao) newSession(table string) *xorm.Session {
	return p.setTable(_db.NewSession(), table)
}

func (p *phoneDao) setTable(sess *xorm.Session, name string) *xorm.Session {
	return sess.Table(p.table + "_" + name)
}

func NewPhoneDao(table string) dao.PhoneDao {
	return &phoneDao{
		table: table,
	}
}
