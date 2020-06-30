package main

import (
	"path/filepath"
	"testing"
)

var (
	templateStoreTestDir string
)

func init() {
	SetIsTestEnvironment(true)
	templateStoreTestDir = filepath.Join(GetExecutableDir(), "testdata", "TemplateStore")
}

func Test_getTemplateStoreForDir_Errors(t *testing.T) {
	// non-existant dir
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "non-existant dir"))
		expectNil(t, ts)
		expectError(t, err)
	}

	// malformed file
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "malformedFile"))
		expectNil(t, ts)
		expectError(t, err)
	}

	// file without description
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "invalidFile"))
		expectNil(t, ts)
		expectError(t, err)
	}

	// files with duplicate IDs
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "duplicateIDs"))
		expectNil(t, ts)
		expectError(t, err)
	}

	// files with cycle inheritance dependencies
	{
		// TODO implement
		/*
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "cyclicTemplateInheritance"))
			expectNil(t, ts)
			expectError(t, err)
		*/
	}
}

func Test_getTemplateStoreForDir_Valid(t *testing.T) {
	// empty dir
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "emptyDir"))
		expectNotNil(t, ts)
		expectNoError(t, err)
		if ts != nil {
			expectEqual(t, len(ts.GetTemplateIDs(true)), 0)
		}
	}

	// valid
	{
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "valid"))
		expectNotNil(t, ts)
		expectNoError(t, err)
		if ts != nil {
			expectEqual(t, len(ts.GetTemplateIDs(true)), 2)
			ct, err := ts.GetTemplate("parent")
			expectNotNil(t, ct)
			expectNoError(t, err)

			ct, err = ts.GetTemplate("child")
			expectNotNil(t, ct)
			expectNoError(t, err)
		}
	}
}

func Test_GetTemplateStore(t *testing.T) {
	SetTestingTemplatesDir(filepath.Join(templateStoreTestDir, "valid"))
	ts, err := GetTemplateStore()

	expectNotNil(t, ts)
	expectNoError(t, err)
	if ts != nil {
		expectTrue(t, len(ts.GetTemplateIDs(false)) > 0)

		ct, err := ts.GetTemplate("parent")
		expectNotNil(t, ct)
		expectNoError(t, err)
	}
}

func Test_TemplateStoreGetTemplate(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "valid"))
	expectNotNil(t, ts)
	expectNoError(t, err)
	if ts != nil {
		ct, err := ts.GetTemplate("parent")
		expectNotNil(t, ct)
		expectNoError(t, err)

		ct, err = ts.GetTemplate("foo")
		expectNil(t, ct)
		expectError(t, err)

		ct, err = ts.GetTemplate("")
		expectNil(t, ct)
		expectError(t, err)
	}

}

func Test_GetTemplate(t *testing.T) {
	SetTestingTemplatesDir(filepath.Join(templateStoreTestDir, "valid"))

	ct, err := GetTemplate("non-existant")
	expectNil(t, ct)
	expectError(t, err)

	ct, err = GetTemplate("parent")
	expectNotNil(t, ct)
	expectNoError(t, err)
}
