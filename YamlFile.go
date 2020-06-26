package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YamlFile represents the structure of a yaml template file
type YamlFile struct {
	Default ContentEntry
	Content []ContentEntry
	//Inherit *string // Name of the template that should be inherited
}

// GetYamlFile reads the yaml file from the provided location.
func GetYamlFile(filename string) (yFile *YamlFile, err error) {
	// TODO print or log reading of yaml file
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
