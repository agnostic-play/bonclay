package httpUtils

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"strings"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplates(path []string) *TemplateRenderer {
	tmpl := template.New("").Funcs(template.FuncMap{
		"trims": strings.TrimSpace,
	})
	for index, val := range path {
		if index == 0 {
			tmpl = template.Must(template.ParseGlob(val))
		}
		tmpl = template.Must(tmpl.ParseGlob(val))
	}

	return &TemplateRenderer{
		templates: tmpl,
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//	Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
