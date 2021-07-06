package config

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// ssr://server:port:protocol:method:obfs:password_base64/?params_base64
type SSRConfig struct {
	Server   string            `json:"server"`
	Port     string            `json:"port"`
	Protocol string            `json:"protocol"`
	Method   string            `json:"method"`
	Obfs     string            `json:"obfs"`
	Password string            `json:"password"`
	Params   map[string]string `json:"params"`
	Raw      string            `json:"raw"`
}

func (s *SSRConfig) Validate() bool {
	// 验证 ip 地址
	ip := net.ParseIP(s.Server)
	if len(ip) == 0 {
		// 检查是否是域名.
		if strings.Contains(s.Server, ".") {
			return true
		}
		return false
	}

	return true
}

type V2rayConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Net  string `json:"net"`
	Raw  string `json:"raw"`
}

func (v V2rayConfig) Validate() bool {
	return v.Host == ""
}

type SSConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Raw  string `json:"raw"`
}

type AllConfig struct {
	SS    []SSConfig    `json:"ss"`
	SSR   []SSRConfig   `json:"ssr"`
	V2ray []V2rayConfig `json:"v2Ray"`
}

type Config struct {
	Domains []*Domain          `yaml:"domains"`
	domains map[string]*Domain `yaml:"domains"`
	Web     Web                `yaml:"web"`
	Updater Updater            `yaml:"updater"`
}

type Domain struct {
	Hosts       []string `yaml:"hosts"`
	GTag        string   `yaml:"gTag"`
	GAd         string   `yaml:"gAd"`
	AdsTxt      string   `yaml:"adsTxt"`
	Theme       string   `yaml:"theme"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	KeyWords    string   `yaml:"keyWords"`
	Html        string   `yaml:"html"`
	Copyright   string   `yaml:"copyright"`
	SiteMap     string   `yaml:"siteMap"`
}

type Web struct {
	Addr           string `yaml:"addr"`
	Port           string `yaml:"port"`
	ForceHTTPS     bool   `yaml:"forceHttps"`
	StaticDir      string `yaml:"staticDir"`
	TemplateDir    string `yaml:"templateDir"`
	EnableCache    bool   `yaml:"enableCache"`
	CacheDirectory string `yaml:"cacheDirectory"`
	CacheDuration  string `yaml:"cacheDuration"`
}

func (c *Config) GetConfigByHost(s string) *Domain {
	return c.domains[s]
}

func (c *Config) MergeMap() {
	c.domains = make(map[string]*Domain)
	for _, d := range c.Domains {
		for _, h := range d.Hosts {
			c.domains[h] = d
		}
	}
}

func GetConfig() *Config {
	return _config
}

var (
	_config *Config
)

func InitConfig() {
	config := os.Getenv("CONFIG")
	if config == "" {
		config = "config.yaml"
	}

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal(data, &_config)
	if err != nil {
		log.Fatal(err)
		return
	}
	_config.MergeMap()
}

func IsProd() bool {
	return gin.Mode() == gin.ReleaseMode
}

type Updater struct {
	Duration        string `json:"duration" yaml:"duration"`
	LastRunTimeFile string `json:"lastRunTimeFile" yaml:"lastRunTimeFile"`
	CachePath       string `json:"cachePath" yaml:"cachePath"`
	Ip2Region       string `json:"ip2Region"  yaml:"ip2Region"`
}
