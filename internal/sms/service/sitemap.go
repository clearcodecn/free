package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/clearcodecn/yunsms/config"
	"github.com/clearcodecn/yunsms/internal/sms/model"
	"github.com/clearcodecn/yunsms/pkg/ext"
	"github.com/sirupsen/logrus"
)

func NewSiteMapService() *SiteMapService {
	return &SiteMapService{
		generator: GetSitemapGenerator(),
	}
}

func (s *SiteMapService) SiteMap(domain *config.Domain, writer io.Writer) {
	fp := s.generator.file(domain.FullHost)
	fi, err := os.OpenFile(fp, os.O_RDONLY, 0777)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read sitemap %v", fp)
		return
	}
	defer fi.Close()
	io.Copy(writer, fi)
}

func (s *SiteMapService) Clean() {
	s.generator.clean()
}

var (
	_siteMapGenerator *siteMapGenerator
	o                 sync.Once
)

func GetSitemapGenerator() *siteMapGenerator {
	o.Do(func() {
		_siteMapGenerator = &siteMapGenerator{
			filePaths: make(map[string]string),
		}
	})
	return _siteMapGenerator
}

type SiteMapService struct {
	generator *siteMapGenerator
}

type siteMapGenerator struct {
	filePaths map[string]string
	mutex     sync.RWMutex
}

func (s *siteMapGenerator) file(host string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if v, ok := s.filePaths[host]; ok {
		return v
	}
	return ""
}

func (s *siteMapGenerator) Start(stopChan chan struct{}) {
	fn := func() {
		cfg := config.GetConfig()
		for _, d := range cfg.Domains {
			fi, _ := os.Stat(d.SiteMap)
			if fi != nil {
				if time.Now().Sub(fi.ModTime()) < 12*time.Hour {
					data, err := ioutil.ReadFile(d.SiteMap)
					if err != nil {
						continue
					}
					s.save(d, string(data))
					continue
				}
			}
			sm := s.generate(d)
			if sm != "" {
				s.save(d, sm)
			}
		}
	}
	fn()
	ticker := time.NewTicker(12 * time.Hour)
	for {
		select {
		case <-stopChan:
			s.clean()
			return
		case <-ticker.C:
			fn()
		}
	}
}

func (s *siteMapGenerator) save(config *config.Domain, sitemap string) {
	os.MkdirAll(filepath.Dir(config.SiteMap), 0777)
	err := ioutil.WriteFile(config.SiteMap, []byte(sitemap), 0777)
	if err != nil {
		logrus.WithError(err).Errorf("failed to write file: %s", config.SiteMap)
		return
	}

	logrus.Infof("generate new sitemap for %s and save at: %s", config.FullHost, config.SiteMap)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.filePaths[config.FullHost] = config.SiteMap
}

func (s *siteMapGenerator) generate(config *config.Domain) string {

	var urls []*Url
	ps, err := s.phoneMessageURLs(config)
	if err != nil {
		return ""
	}
	urls = append(urls, ps...)
	urls = append(urls, s.commonURLs()...)
	urls = append(urls, s.customURLs(config)...)

	data, _ := xml.Marshal(urls)

	siteMap := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset
    xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9
       http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
>
%s
</urlset>
`, string(data))

	return siteMap
}

type Url struct {
	Loc        string   `xml:"loc"`        // url 地址
	Priority   float32  `xml:"priority"`   // 权重
	Lastmod    string   `xml:"lastmod"`    // y-m-d
	Changefreq string   `xml:"changefreq"` // daily
	XMLName    xml.Name `xml:"url"`
}

func (s *siteMapGenerator) phoneMessageURLs(config *config.Domain) ([]*Url, error) {
	page := ext.Pager{
		Page:     1,
		PageSize: 9 * 15,
	}
	svc := NewPhoneService(config.DB)
	msvc := NewMessageService(config.DB)

	list, _, err := svc.ListPhone(context.Background(), page, "")
	if err != nil {
		logrus.WithError(err).Error("phoneXML: failed to listPhone")
		return nil, err
	}
	var urls []*Url

	for _, p := range list {
		count, err := msvc.dao.CountByPhone(context.Background(), p.PhoneNo)
		if err != nil {
			logrus.WithError(err).Error("phoneXML: failed to CountByPhone")
			return nil, err
		}
		urls = append(urls, s.msgUrls(config, p, count)...)
	}

	return urls, nil
}

func getURL(host string, s string) string {
	return host + s
}

func (s *siteMapGenerator) msgUrls(domain *config.Domain, phone *model.Phone, count int64) []*Url {
	var urls []*Url
	var p = 0
	for i := 0; i < int(count); i = i + 15 {
		p++
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/messages/%s?page=%d", phone.Id, p)),
			Priority:   0.8,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
	}

	return urls
}

func (s *siteMapGenerator) commonURLs() []*Url {
	var urls []*Url
	urls = append(urls, &Url{
		Loc:        "/sites/policies",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})
	urls = append(urls, &Url{
		Loc:        "/sites/community",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})
	urls = append(urls, &Url{
		Loc:        "/sites/help",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})
	urls = append(urls, &Url{
		Loc:        "/sites/feedback",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})
	urls = append(urls, &Url{
		Loc:        "/",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})
	urls = append(urls, &Url{
		Loc:        "/site",
		Priority:   0.3,
		Lastmod:    time.Now().Format(`2006-01-02`),
		Changefreq: "daily",
	})

	return urls
}

func (s *siteMapGenerator) customURLs(domain *config.Domain) []*Url {
	var urls []*Url
	if domain.Theme == "usa" {
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/japan")),
			Priority:   0.3,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/russia")),
			Priority:   0.3,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/usa")),
			Priority:   0.3,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/mexico")),
			Priority:   0.3,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
		urls = append(urls, &Url{
			Loc:        getURL(domain.FullHost, fmt.Sprintf("/brazil")),
			Priority:   0.3,
			Lastmod:    time.Now().Format(`2006-01-02`),
			Changefreq: "daily",
		})
	}

	return urls
}

func (s *siteMapGenerator) clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k := range s.filePaths {
		os.RemoveAll(k)
		delete(s.filePaths, k)
	}
}
