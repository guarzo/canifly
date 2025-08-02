package main

import (
	_ "embed"
	"strings"
)

//go:embed version
var versionFile string

// Version is the current application version, loaded from the version file
var Version = strings.TrimSpace(versionFile)