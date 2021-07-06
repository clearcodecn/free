package ext

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin/render"
)

const defaultTheme = "default"

type Render struct {
	o         sync.Once
	templates map[string]*template.Template
	fm        template.FuncMap
	IsProd    bool
	Root      string
}

func (r *Render) Instance(s string, i interface{}) render.Render {
	theme := strings.Split(s, "/")[0]
	if theme == "" {
		theme = defaultTheme
	}

	withoutTheme := strings.TrimPrefix(s, theme+"/")

	var tpl *template.Template
	var err error

	if !r.IsProd {
		t := template.New("").Funcs(r.fm)
		tpl, err = t.ParseGlob(fmt.Sprintf("%s/%s/*.gohtml", r.Root, theme))
	} else {
		r.lazyInit()
		tpl = r.templates[theme]
	}

	if err != nil {
		panic(err)
	}
	return render.HTML{
		Template: tpl,
		Name:     withoutTheme,
		Data:     i,
	}
}

func (r *Render) lazyInit() {
	r.o.Do(func() {
		r.templates = make(map[string]*template.Template)
		for _, theme := range []string{"default"} {
			t := template.New("").Funcs(r.fm)
			tpl, err := t.ParseGlob(fmt.Sprintf("%s/%s/*.gohtml", r.Root, theme))
			if err != nil {
				panic(err)
			}
			r.templates[theme] = tpl
		}
	})
}

func NewRender(root string, isProd bool) render.HTMLRender {
	return &Render{
		fm:     fm,
		Root:   root,
		IsProd: isProd,
	}
}

var fm = template.FuncMap{
	"moment": func(t time.Time) string {
		n := time.Now().Sub(t)
		switch {
		case n.Seconds() < 60:
			return `刚刚`
		case n.Minutes() < 60:
			return fmt.Sprintf("%d分钟前", int(n.Minutes()))
		case n.Hours() < 24:
			return fmt.Sprintf("%d小时前", int(n.Hours()))
		default:
			d := int(n.Hours()) / 24
			return fmt.Sprintf("%d天前", d)
		}
	},
	"momenten": func(t time.Time) string {
		n := time.Now().Sub(t)
		switch {
		case n.Seconds() < 60:
			return `Just Now`
		case n.Minutes() < 60:
			return fmt.Sprintf("%d mins ago", int(n.Minutes()))
		case n.Hours() < 24:
			return fmt.Sprintf("%d hours ago", int(n.Hours()))
		default:
			d := int(n.Hours()) / 24
			return fmt.Sprintf("%d days ago", d)
		}
	},
	"timeFormat": func(t time.Time) string {
		return t.Format(`2006-01-02 15:04`)
	},
	"html": func(s string) template.HTML {
		return template.HTML(s)
	},
	"hideCode": func(s string) template.HTML {
		x := re.ReplaceAllString(s, `****`)
		x = re2.ReplaceAllString(x, "【****】")
		return template.HTML(x)
	},
	"default": func(a, b interface{}) interface{} {
		switch v := b.(type) {
		case string:
			if v == "" {
				return a
			}
			return b
		default:
			return b
		}
	},
}

var re = regexp.MustCompile(`\d{4,6}`)
var re2 = regexp.MustCompile(`【.*?】`)
