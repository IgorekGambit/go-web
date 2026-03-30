package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

const layoutFile = "layouts/base.html"

var (
	htmlTplCache sync.Map // page name -> *template.Template
)

// RenderHTML выставляет заголовки для HTML-ответа, подключает layout (base) и страницу pages/{page}.html
// (в ней должен быть {{define "content"}}…{{end}}), рендерит именованный шаблон "base".
func RenderHTML(w http.ResponseWriter, page string, data any) error {
	if err := validatePageName(page); err != nil {
		return err
	}

	tmpl, err := templateForPage(page)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return err
	}

	setHTMLHeaders(w)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	return err
}

func setHTMLHeaders(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
}

func templateForPage(page string) (*template.Template, error) {
	if t, ok := htmlTplCache.Load(page); ok {
		return t.(*template.Template), nil
	}

	contentFile := fmt.Sprintf("pages/%s.html", page)
	t, err := template.ParseFS(Files, layoutFile, contentFile)
	if err != nil {
		return nil, err
	}

	htmlTplCache.Store(page, t)
	return t, nil
}

func validatePageName(page string) error {
	if page == "" {
		return errors.New("empty page name")
	}
	if strings.ContainsAny(page, `/\`) {
		return errors.New("invalid page name")
	}
	if strings.Contains(page, "..") {
		return errors.New("invalid page name")
	}
	return nil
}
