package main

import "fmt"

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	name  string
	yFile *YamlFile
}

// NewChronicleTemplate converts a YamlFile into a ChronicleTemplate. It returns
// an error if the YamlFile cannot be converted to a ChronicleTemplate, e.g. because
// it is missing required entries.
func NewChronicleTemplate(yFile *YamlFile) (ct *ChronicleTemplate, err error) {
	if yFile == nil {
		return nil, fmt.Errorf("Provided YamlFile ptr was nil")
	}

	ct = new(ChronicleTemplate)

	if !IsSet(yFile.ID) {
		Assert(IsSet(yFile.fileName), "YamlFile filename should always be present")
		return nil, fmt.Errorf("Template file '%v' does not contain an ID", yFile.fileName)
	}
	ct.name = yFile.ID
	ct.yFile = yFile

	// applying default values
	for key, value := range yFile.Content {
		value.applyDefaults(yFile.Default)
		yFile.Content[key] = value
	}

	// TODO check wheth yamlFile is valid, i.e. contains correct content

	return ct, nil
}

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleTemplate
func (ct *ChronicleTemplate) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = ct.yFile.Content[key]
	return
}

// Name returns the name of the Chronicle Template
func (ct *ChronicleTemplate) Name() string {
	return ct.name
}

// Filename returns the file name of the Chronicle Template
func (ct *ChronicleTemplate) Filename() string {
	return ct.yFile.fileName
}
