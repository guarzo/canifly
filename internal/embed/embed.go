package embed

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed static/*
var StaticFiles embed.FS

var StaticFilesSub fs.FS

//go:embed templates/*.tmpl
var templatesFS embed.FS

var Templates *template.Template

func LoadStatic() error {
	var err error
	Templates, err = template.ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		return err
	}
	// Create a sub-filesystem starting at "static"
	StaticFilesSub, err = fs.Sub(StaticFiles, "static")
	if err != nil {
		return err
	}

	return nil
}
