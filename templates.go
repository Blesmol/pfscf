package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

// ContentEntry is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentEntry struct {
	Type     string  // the type which this entry represents
	ID       string  // the ID or name of that concrete content entry
	Desc     string  // Description of this parameter
	X1, Y1   float64 // first set of coordinates
	X2, Y2   float64 // second set of coordinates
	Font     string  // the name of the font (if any) that should be used to display the content
	Fontsize float64 // size of the font in points
	Align    string  // Alignment of the content: L/C/R + T/M/B
	//Flags    []string
}

// YamlFile represents the structure of a yaml template file
type YamlFile struct {
	Default *ContentEntry
	Content *[]ContentEntry
	//Inherit *string // Name of the template that should be inherited
}

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	name    string
	content map[string]ContentEntry
}

// NewChronicleTemplate returns a new ChronicleTemplate object
func NewChronicleTemplate(name string) (c *ChronicleTemplate) {
	c = new(ChronicleTemplate)
	c.name = name
	c.content = make(map[string]ContentEntry)
	return c
}

// TODO #6 generic function for checking required fields in struct
// also output warnings for non-required fields

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(filename string) (yFile *YamlFile, err error) {
	// TODO print or log reading of yaml file
	yFile = new(YamlFile)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileData, yFile)
	if err != nil {
		log.Fatalf("Error parsing yaml file %v: %v\n", filename, err)
		return nil, err
	}

	return yFile, nil
}

// GetTemplateByName returns the template object for the given name, or nil and
// an error object if no template with that name could be found. The template
// name is case-insensitive.
func GetTemplateByName(tmplName string) (yFile *YamlFile, err error) {
	// Keep it simple for the moment. Search in 'templates' subdir
	// for a file with cfgName as basename and 'yml' as file extension

	tmplBaseFilename := strings.ToLower(tmplName) + ".yml"
	tmplFilename := filepath.Join(GetExecutableDir(), "templates", tmplBaseFilename)

	yFile, err = GetYamlFile(tmplFilename)

	return yFile, err
}

// GetChronicleTemplate extracts, processes, and prepares the template
// information from a YamlFile object and puts it into a form
// that can be worked with.
func (yFile *YamlFile) GetChronicleTemplate() (cTmpl *ChronicleTemplate) {
	cTmpl = NewChronicleTemplate("pfs2") // TODO remove hardcoded name

	// add content entries from yamlFile with name mapping into chronicleTemplate
	for _, val := range *yFile.Content {
		Assert(val.ID != "", "No ID provided!")
		id := val.ID
		if _, exists := cTmpl.content[id]; !exists {
			cTmpl.content[id] = val
		} else {
			panic("Duplicate ID found: " + id)
		}
	}

	return cTmpl
}

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleTemplate
func (cTmpl *ChronicleTemplate) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = cTmpl.content[key]
	return
}
