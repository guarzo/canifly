package embed

import (
	"embed"
	"html/template"
)

//go:embed static/plans/* static/invTypes.csv static/.env static/*.png static/*.css static/*.js
var StaticFiles embed.FS

//go:embed templates/*.tmpl
var templatesFS embed.FS

var Templates *template.Template

func LoadTemplates() error {
	var err error
	Templates, err = template.ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		return err
	}
	return nil
}
