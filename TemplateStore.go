package main

// TemplateStore stores multiple ChronicleTemplates and provides means
// to retrieve them by name.
type TemplateStore struct {
	content map[string]*ChronicleTemplate // Store as ptrs so that it is easier to modify them do things like aliasing
}

// GetTemplateStore returns a template store that is already filled with all templates
// contained in the main template directory. If some error showed up during reading and
// parsing files, resolving dependencies etc, then nil is returned together with an error.
func GetTemplateStore() (ts *TemplateStore, err error) {
	return getTemplateStoreForDir(GetTemplatesDir())
}

// getTemplateStoreForDir takes a directory and returns a template store
// for all entries in that directory, including its subdirectories
func getTemplateStoreForDir(dir string) (ts *TemplateStore, err error) {
	return ts, nil
}
