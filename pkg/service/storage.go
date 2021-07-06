package service

import (
	"context"
	"encoding/json"
	"free/internal/sms/model"
	"free/pkg/ext"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Service interface {
	List(ctx context.Context, pager ext.Pager) ([]*model.Link, error)
	Store(ctx context.Context, links []*model.Link) error
}

type fileService struct {
	mutex    sync.RWMutex
	fileName string

	expireTime    time.Time
	cacheDuration time.Duration
	list          []*model.Link
}

func (f *fileService) List(ctx context.Context, pager ext.Pager) ([]*model.Link, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	if time.Now().Before(f.expireTime) {
		// hit cache
		list := copyLinks(f.list)
		return listPage(list, pager)

	} else {
		// 读取文件
		data, err := ioutil.ReadFile(f.fileName)
		if err != nil {
			return nil, err
		}
		var list []*model.Link
		err = json.Unmarshal(data, &list)
		if err != nil {
			return nil, err
		}
		f.mutex.RUnlock()

		f.mutex.Lock()
		f.list = list
		f.expireTime = time.Now().Add(f.cacheDuration)
		f.mutex.Unlock()

		f.mutex.RLock()

		list = copyLinks(f.list)
		return listPage(list, pager)
	}
}

func listPage(list []*model.Link, pager ext.Pager) ([]*model.Link, error) {
	if pager.PageSize == -1 {
		return list, nil
	}
	left := pager.PageSize * (pager.Page - 1)
	right := pager.PageSize + left
	if right > len(list) {
		right = len(list)
	}
	if left < 0 {
		left = 0
	}
	res := list[left:right]
	return res, nil
}

func (f *fileService) Store(ctx context.Context, links []*model.Link) error {
	existing, _ := f.List(ctx, ext.Pager{
		Page:     0,
		PageSize: -1,
	})

	all := merge(existing, links)
	data, err := json.Marshal(all)
	if err != nil {
		return err
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	os.MkdirAll(filepath.Dir(f.fileName), 0777)
	err = ioutil.WriteFile(f.fileName, data, 0777)
	f.expireTime = time.Now()

	logrus.Infof("store data: %d", len(all))

	return err
}

// merge 将 b 往 a 合并.
func merge(a, b []*model.Link) []*model.Link {
	var (
		dict = make(map[string]*model.Link)
	)
	for _, v := range a {
		dict[v.ID] = v
	}
	for _, v := range b {
		dict[v.ID] = v
	}

	var res []*model.Link
	for _, v := range dict {
		res = append(res, v)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreateTime.After(res[j].CreateTime)
	})
	return res
}

func copyLinks(list []*model.Link) []*model.Link {
	var data []*model.Link
	copy(data, list[:])
	return data
}

var (
	_service Service
)

func InitService(filename string, cacheDuration time.Duration) Service {
	_service = &fileService{
		mutex:         sync.RWMutex{},
		fileName:      filename,
		expireTime:    time.Time{},
		cacheDuration: cacheDuration,
		list:          nil,
	}
	return _service
}

func GetService() Service {
	return _service
}
