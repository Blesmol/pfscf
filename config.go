package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

// Config represents the structure of the config yaml file
type Config struct {
	Defaults ConfigDefaults
	Content  *[]ContentEntry
	Inherit  string // Name of the config that should be inherited
}

// ContentEntry is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentEntry struct {
	Type      string  // the type which this entry represents
	ID        string  // the ID or name of that concrete content entry
	Desc      string  // Description of this parameter
	X1, Y1    float64 // first set of coordinates
	X2, Y2    float64 // second set of coordinates
	Font      string  // the name of the font (if any) that should be used to display the content
	Fontsize  float64 // size of the font in points
	Alignment string
	Default   string
	Flags     []string
}

// ConfigDefaults represents all settings for which a default value can be set.
type ConfigDefaults struct {
	Font     string
	Fontsize float64
}

// TODO #6 generic function for checking required fields in struct
// also output warnings for non-required fields

// GetConfigFromFile reads the config file from the provided location.
func GetConfigFromFile(filename string) (c *Config, err error) {
	// print or log reading of config file
	c = new(Config)

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

// GetGlobalConfig reads the configuration from the global config file, which
// should be located in the same directory as the binary
func GetGlobalConfig() (c *Config) {
	c, err := GetConfigFromFile(filepath.Join(GetExecutableDir(), "pfsct.yml"))
	AssertNoError(err)
	return c
}

// GetConfigByName returns the config object for the given name, or nil and
// an error object if no config with that name could be found. The config
// name is case-insensitive.
func GetConfigByName(cfgName string) (c *Config, err error) {
	// Keep it simple for the moment. Search in 'config' subdir
	// for a file with cfgName as basename and 'yml' as file extension

	cfgBaseFilename := strings.ToLower(cfgName) + ".yml"
	cfgFilename := filepath.Join(GetExecutableDir(), "config", cfgBaseFilename)

	c, err = GetConfigFromFile(cfgFilename)

	return c, err
}
