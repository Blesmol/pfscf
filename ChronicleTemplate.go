package main

// ChronicleTemplate represents a template configuration for chronicles. It contains
// information on what to put where.
type ChronicleTemplate struct {
	name    string
	content map[string]ContentEntry
	// TODO Perhaps will include the yaml file instead of the content map?
}

// NewChronicleTemplate returns a new ChronicleTemplate object
func NewChronicleTemplate(name string) (c *ChronicleTemplate) {
	c = new(ChronicleTemplate)
	c.name = name
	c.content = make(map[string]ContentEntry)
	return c
}

// GetContent returns the ContentEntry matching the provided key
// from the current ChronicleTemplate
func (cTmpl *ChronicleTemplate) GetContent(key string) (ce ContentEntry, exists bool) {
	ce, exists = cTmpl.content[key]
	return
}
