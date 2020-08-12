package yaml

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	templateFilePattern = ".+\\.yml$"
)

// GetYamlFilenamesFromDir takes a directory name as input and returns a list of names
// of all yaml files within that dir and its subdirectories. All returned paths are
// prefixed with the provided path argument.
func GetYamlFilenamesFromDir(dirName string) (yamlFilenames []string, err error) {
	tmplFileRegex := regexp.MustCompile(templateFilePattern)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return yamlFilenames, err
	}

	for _, file := range files {
		if file.IsDir() {
			tmplFilesInSubDir, err := GetYamlFilenamesFromDir(filepath.Join(dirName, file.Name()))
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

// ReadYamlFile reads the yaml file from the provided location and stores the
// data into the provided object.
func ReadYamlFile(filename string, ct interface{}) (err error) {
	// TODO add assertion that interface is a ptr
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading file '%v': %v", filename, err)
	}

	err = yaml.Unmarshal(fileData, ct)
	if err != nil {
		return fmt.Errorf("Parsing file '%v': %v", filename, err)
	}

	return nil
}
