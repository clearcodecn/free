package ext

import (
	"bytes"
	"errors"
	"free/assets/static"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

func StaticServer(e *gin.Engine, staticDir string) {
	e.StaticFS("/static", &staticFs{})
	e.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File(filepath.Join(staticDir, "favicon.ico"))
	})
}

type staticFs struct{}

func (s *staticFs) Open(name string) (http.File, error) {
	name = filepath.Join("static", name)
	r := pool.New().(*bytes.Reader)
	data, err := static.Asset(name)
	if err != nil {
		return nil, errors.New("404 not found")
	}
	r.Reset(data)
	return &staticFile{name: name, Reader: r}, nil
}

type staticFile struct {
	name string
	*bytes.Reader
}

func (s *staticFile) Close() error {
	release(s.Reader)
	s.Reader = nil
	return nil
}

func (s staticFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not allowed")
}

func (s staticFile) Stat() (os.FileInfo, error) {
	return static.AssetInfo(s.name)
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewReader(nil)
	},
}

func release(r *bytes.Reader) {
	r.Reset(nil)
	pool.Put(r)
}
