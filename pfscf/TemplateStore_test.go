package main

import (
	"path/filepath"
	"testing"

	test "github.com/Blesmol/pfscf/pfscf/testutils"
	util "github.com/Blesmol/pfscf/pfscf/utils"
)

var (
	templateStoreTestDir        string
	templateStoreInheritTestDir string
)

func init() {
	util.SetIsTestEnvironment(true)
	templateStoreTestDir = filepath.Join(util.GetExecutableDir(), "testdata", "TemplateStore")
	templateStoreInheritTestDir = filepath.Join(templateStoreTestDir, "inheritance")
}

func TestGetTemplateStoreForDir(t *testing.T) {

	t.Run("errors", func(t *testing.T) {
		t.Run("non-existant dir", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "non-existant dir"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("malformed file", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "malformedFile"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("file without description", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "invalidFile"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("files with duplicate IDs", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "duplicateIDs"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("invalid presets", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "invalidPresets"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("invalid content", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "invalidContent"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})
	})

	t.Run("valid", func(t *testing.T) {
		t.Run("empty dir", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "emptyDir"))
			test.ExpectNotNil(t, ts)
			test.ExpectNoError(t, err)
			if ts != nil {
				test.ExpectEqual(t, len(ts.GetTemplateIDs(true)), 0)
			}
		})

		t.Run("valid", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "valid"))
			test.ExpectNotNil(t, ts)
			test.ExpectNoError(t, err)
			if ts != nil {
				test.ExpectEqual(t, len(ts.GetTemplateIDs(true)), 2)
				ct, err := ts.GetTemplate("parent")
				test.ExpectNotNil(t, ct)
				test.ExpectNoError(t, err)

				ct, err = ts.GetTemplate("child")
				test.ExpectNotNil(t, ct)
				test.ExpectNoError(t, err)
			}
		})
	})

	t.Run("inheritance", func(t *testing.T) {
		t.Run("basic validity", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "basicValid"))
			test.ExpectNotNil(t, ts)
			test.ExpectNoError(t, err)
			if ts == nil {
				return
			}

			test.ExpectEqual(t, len(ts.GetTemplateIDs(false)), 6)
			t.Run("id_a", func(t *testing.T) {
				ctA, err := ts.GetTemplate("id_a")
				test.ExpectNoError(t, err)
				test.ExpectNotSet(t, ctA.Inherit())

				contentIDList := ctA.GetContentIDs(false)
				test.ExpectEqual(t, len(contentIDList), 1)
				test.ExpectEqual(t, contentIDList[0], "content_a")

				presetIDList := ctA.presets.GetIDs()
				test.ExpectEqual(t, len(presetIDList), 1)
				test.ExpectEqual(t, presetIDList[0], "preset_a")
			})

			t.Run("id_c", func(t *testing.T) {
				ctC, err := ts.GetTemplate("id_c")
				test.ExpectNoError(t, err)
				test.ExpectEqual(t, ctC.Inherit(), "id_a")

				contentIDList := ctC.GetContentIDs(false)
				test.ExpectEqual(t, len(contentIDList), 3)
				test.ExpectEqual(t, contentIDList[0], "content_a")
				test.ExpectEqual(t, contentIDList[1], "content_c1")
				test.ExpectEqual(t, contentIDList[2], "content_c2")

				presetIDList := ctC.presets.GetIDs()
				test.ExpectEqual(t, len(presetIDList), 3)
				test.ExpectEqual(t, presetIDList[0], "preset_a")
				test.ExpectEqual(t, presetIDList[1], "preset_c1")
				test.ExpectEqual(t, presetIDList[2], "preset_c2")
			})

			t.Run("id_d", func(t *testing.T) {
				ctD, err := ts.GetTemplate("id_d")
				test.ExpectNoError(t, err)
				test.ExpectEqual(t, ctD.Inherit(), "id_c")

				contentIDList := ctD.GetContentIDs(false)
				test.ExpectEqual(t, len(contentIDList), 4)
				test.ExpectEqual(t, contentIDList[0], "content_a")
				test.ExpectEqual(t, contentIDList[1], "content_c1")
				test.ExpectEqual(t, contentIDList[2], "content_c2")
				test.ExpectEqual(t, contentIDList[3], "content_d")

				presetIDList := ctD.presets.GetIDs()
				test.ExpectEqual(t, len(presetIDList), 4)
				test.ExpectEqual(t, presetIDList[0], "preset_a")
				test.ExpectEqual(t, presetIDList[1], "preset_c1")
				test.ExpectEqual(t, presetIDList[2], "preset_c2")
				test.ExpectEqual(t, presetIDList[3], "preset_d")
			})
		})

		t.Run("cyclic inheritance across multiple files", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "cyclicDep"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("cyclic inheritance from self", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "fromSelf"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})

		t.Run("invalid dependency", func(t *testing.T) {
			ts, err := getTemplateStoreForDir(filepath.Join(templateStoreInheritTestDir, "fromNonExisting"))
			test.ExpectNil(t, ts)
			test.ExpectError(t, err)
		})
	})
}

func TestGetTemplateStore(t *testing.T) {
	SetTestingTemplatesDir(filepath.Join(templateStoreTestDir, "valid"))
	ts, err := GetTemplateStore()

	test.ExpectNotNil(t, ts)
	test.ExpectNoError(t, err)
	if ts != nil {
		test.ExpectTrue(t, len(ts.GetTemplateIDs(false)) > 0)

		ct, err := ts.GetTemplate("parent")
		test.ExpectNotNil(t, ct)
		test.ExpectNoError(t, err)
	}
}

func TestTemplateStore_GetTemplate(t *testing.T) {
	ts, err := getTemplateStoreForDir(filepath.Join(templateStoreTestDir, "valid"))
	test.ExpectNotNil(t, ts)
	test.ExpectNoError(t, err)
	if ts != nil {
		ct, err := ts.GetTemplate("parent")
		test.ExpectNotNil(t, ct)
		test.ExpectNoError(t, err)

		ct, err = ts.GetTemplate("foo")
		test.ExpectNil(t, ct)
		test.ExpectError(t, err)

		ct, err = ts.GetTemplate("")
		test.ExpectNil(t, ct)
		test.ExpectError(t, err)
	}

}

func TestGetTemplate(t *testing.T) {
	SetTestingTemplatesDir(filepath.Join(templateStoreTestDir, "valid"))

	ct, err := GetTemplate("non-existant")
	test.ExpectNil(t, ct)
	test.ExpectError(t, err)

	ct, err = GetTemplate("parent")
	test.ExpectNotNil(t, ct)
	test.ExpectNoError(t, err)
}
