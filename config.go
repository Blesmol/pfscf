package main

import (
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

// Config represents the structure of the config yaml file
type Config struct {
	Defaults ConfigDefaults
	Content  *[]ContentEntry
}

// ContentEntry is a generic struct with lots of fields to fit all
// supported types of Content. Each type will only check its required
// fields. So basically only field "Type" always has to be provided,
// everything else depends on the concrete type.
type ContentEntry struct {
	Type      string  // the type which this entry represents
	ID        string  // the ID or name of that concrete content entry
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
func GetConfigFromFile(filename string) (c *Config) {
	// print or log reading of config file
	c = new(Config)

	fileData, err := ioutil.ReadFile(filename)
	AssertNoError(err)

	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		log.Fatalf("Error parsing config file %v: %v\n", filename, err)
	}

	return c
}

// GetGlobalConfig reads the configuration from the global config file, which
// should be located in the same directory as the binary
func GetGlobalConfig() (c *Config) {
	c = GetConfigFromFile("pfsct.yml")
	//c = GetConfigFromFile(filepath.Join(GetExecutableDir(), "test.yml")
	// TODO check environment for dir info
	return c
}
