package template

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/reyhardy/exs-prescription-app/internal/ctx"
)

type PageName string

const (
	Login              PageName = "login"
	UserProfile        PageName = "profile"
	PatientsList       PageName = "patients_list"
	PatientView        PageName = "patient_view"
	PrescriptionDetail PageName = "prescription_detail"
)

var (
	ErrPageNotExists = errors.New("page not exists")
)

//go:embed all:layout all:page
var PageView embed.FS

type Renderer interface {
	Render(context.Context, io.Writer, PageName, any) error
}

type Manager struct {
	templates map[PageName]*template.Template
}

func NewManager() (*Manager, error) {
	layouts, err := fs.Glob(PageView, "layout/*.gohtml")
	if err != nil {
		return nil, err
	}
	pages, err := fs.Glob(PageView, "page/*.gohtml")
	if err != nil {
		return nil, err
	}
	tMap, err := parseTemplates(pages, layouts...)
	if err != nil {
		return nil, err
	}
	return &Manager{
		templates: tMap,
	}, nil
}

func parseTemplates(pages []string, layouts ...string) (map[PageName]*template.Template, error) {
	templatesMap := make(map[PageName]*template.Template)
	for _, page := range pages {
		files := append(
			append([]string{}, layouts...),
			page,
		)
		tmpl, err := template.ParseFS(PageView, files...)
		if err != nil {
			return nil, err
		}
		page = strings.TrimSuffix(path.Base(page), ".gohtml")
		templatesMap[PageName(page)] = tmpl
	}
	fmt.Printf("templates map: %#v\n", templatesMap)
	return templatesMap, nil
}

func (m *Manager) Render(context context.Context, w io.Writer, name PageName, data any) error {
	t, ok := m.templates[name]
	if !ok {
		return ErrPageNotExists
	}
	uCtx := ctx.GetContext(context)
	pageData := struct {
		Data any
		User *ctx.UserCtx
	}{
		Data: data,
		User: uCtx,
	}
	buf := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buf, "base.gohtml", pageData); err != nil {
		return err
	}
	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}
