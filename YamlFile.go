package main

import (
	"fmt"
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
	Inherit     string                  // ID of the template that should be inherited
	Default     ContentEntry            // default values for the Content entries
	Presets     map[string]ContentEntry // Named preset sections
	Content     map[string]ContentEntry // The Content.
}

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(filename string) (yFile *YamlFile, err error) {
	yFile = new(YamlFile)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Reading file '%v': %w", filename, err)
	}

	err = yaml.UnmarshalStrict(fileData, yFile)
	if err != nil {
		return nil, fmt.Errorf("Parsing file '%v': %w", filename, err)
	}

	// set content id inside presets entries
	for id, entry := range yFile.Presets {
		Assert(!IsSet(entry.id), "ContentEnty id should not be already set")
		entry.id = id
		yFile.Presets[id] = entry
	}

	// set content id inside content entries
	for id, entry := range yFile.Content {
		Assert(!IsSet(entry.id), "ContentEnty id should not be already set")
		entry.id = id
		yFile.Content[id] = entry
	}

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
			filename := filepath.Join(dirName, file.Name())
			yamlFilenames = append(yamlFilenames, filename)
		}
	}
	return yamlFilenames, nil
}

// GetTemplateFilesFromDir takes a directory name as input and returns a list of
// YamlFile objects that hold the contents of all yaml files contained in that
// directory and its subdirectories.
func GetTemplateFilesFromDir(dirName string) (yamlFiles map[string]*YamlFile, err error) {
	fileList, err := GetTemplateFilenamesFromDir(dirName)
	if err != nil {
		return nil, err
	}

	yamlFiles = make(map[string]*YamlFile, len(fileList))
	for _, filename := range fileList {
		yFile, err := GetYamlFile(filename)
		if err != nil {
			return nil, err
		}

		yamlFiles[filename] = yFile
	}

	return yamlFiles, nil
}
