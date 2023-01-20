package httpUtils

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

var tmpl embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplates(path string) *TemplateRenderer {

	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./resources/views/*.gohtml")),
		//templates: template.Must(template.ParseFS(tmpl, "./resources/views/*.gohtml")),
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	//
	//	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	//
	return t.templates.ExecuteTemplate(w, name, data)
}
