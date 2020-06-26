package main

import (
	"path/filepath"
)

const (
	templateDir = "templates"
)

// GetTemplatesDir returns the path of the top directory where the template files are stored
func GetTemplatesDir() (dir string) {
	return filepath.Join(GetExecutableDir(), templateDir)
}
