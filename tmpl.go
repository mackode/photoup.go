package main

import (
	"io"
	tmpl "text/template"
	"time"
)

type Photo struct {
	Path string
	Thumb string
}

type TmplData struct {
	CGI string
	TmplEngine *tmpl.Template
	Date string
	OgDesc string
	OgImage string
	Photos []Photo
	Link string
}

func NewTmpl() *TmplData {
	return &TmplData{}
}

func (td *TmplData) Init() error {
	td.Date = time.Now().Format("2006-01-02")
	td.OgDesc = "A collection of photos"
	td.CGI = "cgi/photoup"
	te, err := tmpl.ParseGlob("tmpl/*.html")
	td.TmplEngine = te
	return err
}

func (td *TmplData) AddPhoto(name string) {
	td.Photos = append(td.Photos, Photo{Path: name, Thumb: thumbName(name)})
}

func (td *TmplData) RenderPage(w io.Writer, tmplName string) error {
	if len(td.Photos) != 0 {
		td.OgImage = td.Photos[0].Thumb
	}

	for _, tpl := range []string{"intro.html", tmplName, "outro.html"} {
		err := td.TmplEngine.ExecuteTemplate(w, tpl, td)
		if err != nil {
			return err
		}
	}

	return nil
}
