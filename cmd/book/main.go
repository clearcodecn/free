package main

import (
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	tpl *template.Template
)

func init() {
	data, err := ioutil.ReadFile("template/render.html")
	if err != nil {
		panic(err)
	}
	tpl, err = template.New("").Parse(string(data))
	if err != nil {
		panic(err)
	}
}

// 将 book 渲染成 html.
func main() {
	err := filepath.Walk("./books", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		sfx := strings.ToLower(filepath.Ext(path))
		if sfx != ".md" {
			return err
		}
		return render(path)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func render(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	output := blackfriday.MarkdownBasic(data)

	var buf bytes.Buffer
	err = tpl.Execute(&buf, template.HTML(string(output)))
	if err != nil {
		return err
	}

	fileName := filepath.Base(path)
	name := strings.TrimRight(fileName, filepath.Ext(path))

	outPath := filepath.Join("static", "article", name+".html")
	os.MkdirAll(filepath.Dir(outPath), 0777)
	err = ioutil.WriteFile(outPath, buf.Bytes(), 0777)
	return err
}
