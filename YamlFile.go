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
	ID          string                  // Name by which this template should be identified
	Description string                  // The description of this template
	Default     ContentEntry            // default values for the Content entries
	Content     map[string]ContentEntry // The Content.
	fileName    string                  // not exported, as this field should not be set via the yaml file
	//Inherit string // Name of the template that should be inherited
}

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(fileName string) (yFile *YamlFile, err error) {
	yFile = new(YamlFile)

	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.UnmarshalStrict(fileData, yFile)
	if err != nil {
		return nil, err
	}

	yFile.fileName = fileName

	return yFile, nil
}

// GetTemplateFilenamesFromDir takes a directory name as input and returns a list of names
// of all template files within that dir and its subdirectories. All returned paths are
// prefixed with the provided path argument.
func GetTemplateFilenamesFromDir(dirName string) (yamlFilenames []string, err error) {
	tmplFileRegex := regexp.MustCompile(templateFilePattern)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return yamlFilenames, err
	}

	for _, file := range files {
		if file.IsDir() {
			tmplFilesInSubDir, err := GetTemplateFilenamesFromDir(filepath.Join(dirName, file.Name()))
			if err != nil {
				return nil, err
			}
			yamlFilenames = append(yamlFilenames, tmplFilesInSubDir...)
		} else if tmplFileRegex.MatchString(strings.ToLower(file.Name())) {
			fileName := filepath.Join(dirName, file.Name())
			yamlFilenames = append(yamlFilenames, fileName)
		}
	}
	return yamlFilenames, nil
}

// GetTemplateFilesFromDir takes a directory name as input and returns a list of
// YamlFile objects that hold the contents of all yaml files contained in that
// directory and its subdirectories.
func GetTemplateFilesFromDir(dirName string) (yamlFiles []*YamlFile, err error) {
	fileList, err := GetTemplateFilenamesFromDir(dirName)
	if err != nil {
		return nil, err
	}

	for _, fileName := range fileList {
		yFile, err := GetYamlFile(fileName)
		if err != nil {
			return nil, err
		}

		yamlFiles = append(yamlFiles, yFile)
	}

	return yamlFiles, nil
}
