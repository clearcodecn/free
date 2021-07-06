package ext

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CacheConfig struct {
	Directory string
	Duration  string
	Enable    bool
}

func Cache(config CacheConfig) gin.HandlerFunc {
	if !config.Enable {
		return gin.HandlerFunc(func(ctx *gin.Context) {
			ctx.Next()
		})
	}

	os.MkdirAll(config.Directory, 0777)
	mu := sync.Mutex{}

	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		duration = time.Second * 30
	}

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.Next()
			return
		}

		if strings.HasPrefix(ctx.Request.RequestURI, "/static") {
			ctx.Next()
			return
		}

		if strings.HasPrefix(ctx.Request.RequestURI, "/favicon.ico") {
			ctx.Next()
			return
		}

		mu.Lock()
		hash := buildHash(ctx.Request)
		path := filepath.Join(config.Directory, hash)
		fi, err := os.Stat(path)

		hit := true
		if err != nil {
			if !os.IsNotExist(err) {
				mu.Unlock()

				ctx.AbortWithError(503, err)
				return
			}
			hit = false
		}

		if fi != nil {
			if time.Now().Sub(fi.ModTime()) > duration {
				hit = false
			} else {
				hit = true
			}
		}

		if hit {
			mu.Unlock()

			ctx.Header("Content-Type", "text/html")
			ctx.File(path)
			ctx.Abort()
			logrus.Infof("hit cache %s %s %s", path, ctx.Request.Host, ctx.Request.URL.String())
			return
		}
		mu.Unlock()

		w := &responseWriter{ResponseWriter: ctx.Writer}
		ctx.Writer = w

		ctx.Next()
		if w.status != 200 {
			w.realWrite()
			return
		}

		mu.Lock()
		ioutil.WriteFile(path, w.buf, 0777)
		mu.Unlock()

		w.realWrite()
	}
}

func buildHash(r *http.Request) string {
	m := md5.New()
	m.Write([]byte(r.Host + r.URL.String()))
	return fmt.Sprintf("%x.html", m.Sum(nil))
}

type responseWriter struct {
	gin.ResponseWriter
	buf    []byte
	status int
}

func (r *responseWriter) Write(b []byte) (int, error) {
	r.buf = append(r.buf, b...)
	return len(b), nil
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseWriter) realWrite() {
	r.ResponseWriter.WriteHeader(r.status)
	r.ResponseWriter.Write(r.buf)
}
