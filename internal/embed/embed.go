package embed

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var StaticFiles embed.FS

var StaticFilesSub fs.FS

//go:embed config/.env
var EnvFiles embed.FS

func LoadStatic() error {
	var err error

	// Create a sub-filesystem starting at "static"
	StaticFilesSub, err = fs.Sub(StaticFiles, "static")
	if err != nil {
		return err
	}

	return nil
}
