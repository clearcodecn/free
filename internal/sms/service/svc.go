package service

import (
	"encoding/json"
	"errors"
	"free/config"
	"free/internal/sms/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"runtime/debug"
	"sync"
	"time"
)

type LinkService struct {
	cache  []*model.Link
	mu     sync.RWMutex
	update chan []*model.Link
}

var (
	_service *LinkService
)

func NewLinkService() *LinkService {
	if _service == nil {
		_service = &LinkService{
			update: make(chan []*model.Link),
		}
		go _service.Do()
	}
	return _service
}

func (l *LinkService) FindAll() []*model.Link {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var links = make([]*model.Link, len(l.cache))
	copy(links, l.cache[:])

	return links
}

func (l *LinkService) FindOne(id string) (*model.Link, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, item := range l.cache {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, errors.New("not found")
}

func (l *LinkService) Do() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	fn := func() {
		file := config.GetConfig().Updater.CachePath
		data, err := ioutil.ReadFile(file)
		if err != nil {
			logrus.WithError(err).Errorf("failed to read file")
			return
		}
		var list []*model.Link
		err = json.Unmarshal(data, &list)
		if err != nil {
			logrus.WithError(err).Errorf("failed to Unmarshal data")
			return
		}
		l.mu.Lock()
		l.cache = list
		l.mu.Unlock()
	}
	fn()

	tk := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-tk.C:
			fn()
		case list := <-l.update:
			l.mu.Lock()
			l.cache = list
			l.mu.Unlock()
		}
	}
}

func (l *LinkService) Update(list []*model.Link) {
	l.update <- list
}
