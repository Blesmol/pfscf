package main

import (
	"path/filepath"
)

// GetTemplatesDir returns the path of the top directory where the template files are stored
func GetTemplatesDir() (dir string) {
	return filepath.Join(GetExecutableDir(), "templates")
}
