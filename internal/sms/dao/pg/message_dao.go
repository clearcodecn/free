package pg

import (
	"context"

	"github.com/clearcodecn/yunsms/config"
	"github.com/clearcodecn/yunsms/internal/sms/dao"
	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
	"github.com/go-xorm/xorm"
)

var (
	_db          *xorm.Engine
	phoneTable   = "phone"
	messageTable = "message"
)

func InitDB() {
	cfg := config.GetConfig().PgConfig
	e, err := xorm.NewEngine("postgres", cfg.DSN)
	if err != nil {
		panic(err)
	}
	_db = e

	//e.ShowSQL(true)
	for _, d := range config.GetConfig().Domains {
		sess := e.NewSession()

		err = sess.Table(d.DB + "_" + phoneTable).Sync2(new(model.Phone))
		if err != nil {
			panic(err)
		}
		err = sess.Table(d.DB + "_" + messageTable).Sync2(new(model.Message))
		if err != nil {
			panic(err)
		}
		sess.Close()
	}
}

func (m *messageDao) newSession(table string) *xorm.Session {
	return m.setTable(_db.NewSession(), table)
}

type messageDao struct {
	table string
}

func (m *messageDao) setTable(sess *xorm.Session, name string) *xorm.Session {
	return sess.Table(m.table + "_" + name)
}

func (m *messageDao) ListMessageByPhone(ctx context.Context, phone string, pager ext.Pager) ([]*model.Message, int64, error) {
	sess := m.newSession(messageTable)
	defer sess.Close()

	sess.Where("phone_no = ?", phone)
	sess.Limit(pager.PageSize, pager.From())
	sess.OrderBy("create_time desc")

	var message []*model.Message
	count, err := sess.FindAndCount(&message)
	if err != nil {
		return nil, 0, err
	}

	return message, count, nil
}

func (m *messageDao) CreateMessage(ctx context.Context, sms *model.Message) error {
	sess := m.newSession(messageTable)
	defer sess.Close()

	_, err := sess.InsertOne(sms)
	if err != nil {
		return err
	}
	return nil
}

func (m *messageDao) BulkCreateMessage(ctx context.Context, sms []*model.Message) error {
	sess := m.newSession(messageTable)
	defer sess.Close()

	for {
		var data []*model.Message
		if len(sms) > 10000 {
			data = sms[:10000]
			sms = sms[10000:]
		} else {
			data = sms
			sms = sms[:0]
		}
		if _, err := m.setTable(sess, messageTable).InsertMulti(data); err != nil {
			return err
		}
		if len(sms) == 0 {
			break
		}
	}
	return nil
}

func (m *messageDao) CountByPhone(ctx context.Context, phone string) (int64, error) {
	sess := m.newSession(messageTable)
	defer sess.Close()

	sess.Where("phone_no = ?", phone)
	return sess.Count(new(model.Message))
}

func NewMessageDao(table string) dao.MessageDao {
	return &messageDao{
		table: table,
	}
}
