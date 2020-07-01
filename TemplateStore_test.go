package main

import (
	"path/filepath"
	"testing"
)

var (
	templateStoreTestDir        string
	templateStoreInheritTestDir string
)

func init() {
	SetIsTestEnvironment(true)
	templateStoreTestDir = filepath.Join(GetExecutableDir(), "testdata", "TemplateStore")
	templateStoreInheritTestDir = filepath.Join(templateStoreTestDir, "inheritance")
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

func Test_getTemplateStoreForDir_Inheritance(t *testing.T) {

	t.Run("basic validity", func(t *testing.T) {
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "basicValid"))
		expectNotNil(t, ts)
		expectNoError(t, err)
		if ts == nil {
			return
		}

		expectEqual(t, len(ts.GetTemplateIDs(false)), 6)
		t.Run("id_a", func(t *testing.T) {
			ctA, err := ts.GetTemplate("id_a")
			expectNoError(t, err)
			expectNotSet(t, ctA.Inherit())

			contentIDList := ctA.GetContentIDs(false)
			expectEqual(t, len(contentIDList), 1)
			expectEqual(t, contentIDList[0], "content_a")

			presetIDList := ctA.GetPresetIDs()
			expectEqual(t, len(presetIDList), 1)
			expectEqual(t, presetIDList[0], "preset_a")
		})

		t.Run("id_c", func(t *testing.T) {
			ctC, err := ts.GetTemplate("id_c")
			expectNoError(t, err)
			expectEqual(t, ctC.Inherit(), "id_a")

			contentIDList := ctC.GetContentIDs(false)
			expectEqual(t, len(contentIDList), 3)
			expectEqual(t, contentIDList[0], "content_a")
			expectEqual(t, contentIDList[1], "content_c1")
			expectEqual(t, contentIDList[2], "content_c2")

			presetIDList := ctC.GetPresetIDs()
			expectEqual(t, len(presetIDList), 3)
			expectEqual(t, presetIDList[0], "preset_a")
			expectEqual(t, presetIDList[1], "preset_c1")
			expectEqual(t, presetIDList[2], "preset_c2")
		})

		t.Run("id_d", func(t *testing.T) {
			ctD, err := ts.GetTemplate("id_d")
			expectNoError(t, err)
			expectEqual(t, ctD.Inherit(), "id_c")

			contentIDList := ctD.GetContentIDs(false)
			expectEqual(t, len(contentIDList), 4)
			expectEqual(t, contentIDList[0], "content_a")
			expectEqual(t, contentIDList[1], "content_c1")
			expectEqual(t, contentIDList[2], "content_c2")
			expectEqual(t, contentIDList[3], "content_d")

			presetIDList := ctD.GetPresetIDs()
			expectEqual(t, len(presetIDList), 4)
			expectEqual(t, presetIDList[0], "preset_a")
			expectEqual(t, presetIDList[1], "preset_c1")
			expectEqual(t, presetIDList[2], "preset_c2")
			expectEqual(t, presetIDList[3], "preset_d")
		})
	})

	t.Run("cyclic inheritance across multiple files", func(t *testing.T) {
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "cyclicDep"))
		expectNil(t, ts)
		expectError(t, err)
	})

	t.Run("cyclic inheritance from self", func(t *testing.T) {
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "fromSelf"))
		expectNil(t, ts)
		expectError(t, err)
	})

	t.Run("invalid dependency", func(t *testing.T) {
		ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "fromNonExisting"))
		expectNil(t, ts)
		expectError(t, err)
	})

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
