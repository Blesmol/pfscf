package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	templateFilePattern = ".+\\.yml"
)

// YamlFile represents the structure of a yaml template file
type YamlFile struct {
	ID          string       // Name by which this template should be identified
	Description string       // The description of this template
	Default     ContentEntry // default values for the Content entries
	Content     []ContentEntry
	//Inherit string // Name of the template that should be inherited
}

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(filename string) (yFile *YamlFile, err error) {
	yFile = new(YamlFile)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.UnmarshalStrict(fileData, yFile)
	if err != nil {
		return nil, err
	}

	return yFile, nil
}

// GetTemplateFilenamesFromDir takes a directory name as input and returns a list of names
// of all template files within that dir and its subdirectories. All returned paths are
// prefixed with the provided path argument.
func GetTemplateFilenamesFromDir(dirName string) (tmplFiles []string, err error) {
	tmplFileRegex := regexp.MustCompile(templateFilePattern)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return tmplFiles, err
	}

	for _, file := range files {
		if file.IsDir() {
			tmplFilesInSubDir, err := GetTemplateFilenamesFromDir(filepath.Join(dirName, file.Name()))
			if err != nil {
				return nil, err
			}
			tmplFiles = append(tmplFiles, tmplFilesInSubDir...)
		} else if tmplFileRegex.MatchString(strings.ToLower(file.Name())) {
			fileName := filepath.Join(dirName, file.Name())
			tmplFiles = append(tmplFiles, fileName)
		}
	}
	return
}
