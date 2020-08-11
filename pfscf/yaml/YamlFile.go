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

// File represents the structure of a yaml template file
type File struct {
	ID          string                 // Name by which this template should be identified
	Description string                 // The description of this template
	Inherit     string                 // ID of the template that should be inherited
	Presets     map[string]ContentData // Named preset sections
	Content     map[string]ContentData // The Content.
}

// ContentData is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentData struct {
	Type     string   // the type which this entry represents
	Desc     string   // Description of this parameter
	X1       float64  `yaml:"x"` // first x coordinate
	Y1       float64  `yaml:"y"` // first y coordinate
	X2, Y2   float64  // second set of coordinates
	XPivot   float64  // pivot point on X axis
	Font     string   // the name of the font (if any) that should be used to display the content
	Fontsize float64  // size of the font in points
	Align    string   // Alignment of the content: L/C/R + T/M/B
	Color    string   // Color code
	Example  string   // Example value to be displayed to users
	Presets  []string // List of presets that should be applied on this ContentData / ContentEntry
	//Flags    *[]string
}

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(filename string) (yFile *File, err error) {
	yFile = new(File)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Reading file '%v': %v", filename, err)
	}

	err = yaml.UnmarshalStrict(fileData, yFile)
	if err != nil {
		return nil, fmt.Errorf("Parsing file '%v': %v", filename, err)
	}

	return yFile, nil
}

// getTemplateFilenamesFromDir takes a directory name as input and returns a list of names
// of all template files within that dir and its subdirectories. All returned paths are
// prefixed with the provided path argument.
func getTemplateFilenamesFromDir(dirName string) (yamlFilenames []string, err error) {
	tmplFileRegex := regexp.MustCompile(templateFilePattern)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return yamlFilenames, err
	}

	for _, file := range files {
		if file.IsDir() {
			tmplFilesInSubDir, err := getTemplateFilenamesFromDir(filepath.Join(dirName, file.Name()))
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
func GetTemplateFilesFromDir(dirName string) (yamlFiles map[string]*File, err error) {
	fileList, err := getTemplateFilenamesFromDir(dirName)
	if err != nil {
		return nil, err
	}

	yamlFiles = make(map[string]*File, len(fileList))
	for _, filename := range fileList {
		yFile, err := GetYamlFile(filename)
		if err != nil {
			return nil, err
		}

		yamlFiles[filename] = yFile
	}

	return yamlFiles, nil
}
