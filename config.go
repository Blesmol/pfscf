package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

// YamlConfig represents the structure of the config yaml file
type YamlConfig struct {
	Default ConfigDefaults
	Content *[]ContentEntry
	Inherit string // Name of the config that should be inherited
}

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
	Default  string
	Flags    []string
}

// ConfigDefaults represents all settings for which a default value can be set.
type ConfigDefaults struct {
	Font string
}

// ChronicleConfig represents a configuration for chronicles. It contains
// information on what to put where.
type ChronicleConfig struct {
	name    string
	content map[string]ContentEntry
}

// NewChronicleConfig returns a new ChronicleConfig object
func NewChronicleConfig(name string) (c *ChronicleConfig) {
	c = new(ChronicleConfig)
	c.name = name
	c.content = make(map[string]ContentEntry)
	return c
}

// TODO #6 generic function for checking required fields in struct
// also output warnings for non-required fields

// GetYamlConfigFromFile reads the config file from the provided location.
func GetYamlConfigFromFile(filename string) (c *YamlConfig, err error) {
	// print or log reading of config file
	c = new(YamlConfig)

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		log.Fatalf("Error parsing config file %v: %v\n", filename, err)
		return nil, err
	}

	return c, nil
}

// GetConfigByName returns the config object for the given name, or nil and
// an error object if no config with that name could be found. The config
// name is case-insensitive.
func GetConfigByName(cfgName string) (c *YamlConfig, err error) {
	// Keep it simple for the moment. Search in 'config' subdir
	// for a file with cfgName as basename and 'yml' as file extension

	cfgBaseFilename := strings.ToLower(cfgName) + ".yml"
	cfgFilename := filepath.Join(GetExecutableDir(), "config", cfgBaseFilename)

	c, err = GetYamlConfigFromFile(cfgFilename)

	return c, err
}

// GetChronicleConfig extracts, processes, and prepares the config
// information from a YamlConfig object and puts it into a form
// that can be worked with.
func (yCfg *YamlConfig) GetChronicleConfig() (cCfg *ChronicleConfig) {
	cCfg = NewChronicleConfig("pfs2") // TODO remove hardcoded name

	// add content entries from yamlConfig with name mapping into chronicleConfig
	for _, val := range *yCfg.Content {
		Assert(val.ID != "", "No ID provided!")
		id := val.ID
		if _, exists := cCfg.content[id]; !exists {
			cCfg.content[id] = val
		} else {
			panic("Duplicate ID found: " + id)
		}
	}

	return cCfg
}

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleConfig
func (cCfg *ChronicleConfig) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = cCfg.content[key]
	return
}
