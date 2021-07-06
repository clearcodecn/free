package updater

import (
	"context"
	"free/config"
	"free/internal/sms/model"
	"free/pkg/service"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	// 时间间隔.
	Duration time.Duration
	// 记录上一次时间的文件, 方便断开了下一次还能接着等待.
	LastRunTimFile string
	// 记录的文件
	StorageFile string
}

func Do(stopChan chan struct{}, config *Config) {

	// 1. 判断一下是否应该立即开始.
	durs, ok := shouldRun(config)
	if ok {
		doOnce(config)
		durs = config.Duration
	}

	var (
		tk = time.NewTicker(durs)
	)
	for {
		select {
		case <-tk.C:
			doOnce(config)
			tk.Stop()

			tk = time.NewTicker(config.Duration)
		case <-stopChan:
			tk.Stop()
			return
		}
	}
}

func storeLastRunTime(config *Config) {
	os.MkdirAll(filepath.Dir(config.StorageFile), 0777)
	t := time.Now().Format(time.RFC3339)
	err := ioutil.WriteFile(config.LastRunTimFile, []byte(t), 0777)
	if err != nil {
		logrus.WithError(err).Errorf("failed to save last run time file")
	}
}

func shouldRun(config *Config) (time.Duration, bool) {
	// 1. 上一次记录的时间.
	data, err := ioutil.ReadFile(config.LastRunTimFile)
	if err != nil {
		return 0, true
	}

	ti, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return 0, true
	}

	durs := time.Now().Sub(ti)
	if durs > config.Duration {
		return 0, true
	}

	return config.Duration - durs, false
}

func allToLinks(configs config.AllConfig) []*model.Link {
	var links []*model.Link
	for _, s := range configs.SS {
		if s.Host == "" {
			continue
		}
		links = append(links, &model.Link{
			ID:         s.Host + s.Port,
			Host:       s.Host,
			Port:       s.Port,
			Raw:        s.Raw,
			Type:       model.LinkTypeSS,
			CreateTime: time.Now(),
		})
	}

	for _, s := range configs.SSR {
		if s.Server == "" {
			continue
		}
		links = append(links, &model.Link{
			ID:         s.Server + s.Port,
			Host:       s.Server,
			Port:       s.Port,
			Raw:        s.Raw,
			Type:       model.LinkTypeSSR,
			CreateTime: time.Now(),
		})
	}

	for _, s := range configs.V2ray {
		if s.Host == "" || s.Port == "" {
			continue
		}
		links = append(links, &model.Link{
			ID:         s.Host + s.Port,
			Host:       s.Host,
			Port:       s.Port,
			Raw:        s.Raw,
			Type:       model.LinkTypeV2ray,
			CreateTime: time.Now(),
		})
	}

	return links
}

func doOnce(config *Config) {
	all := Update()
	links := allToLinks(all)
	if err := service.GetService().Store(context.Background(), links); err != nil {
		logrus.WithError(err).Errorf("failed to store data")
	} else {
		storeLastRunTime(config)
	}
}
