package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

var (
	isTestEnvironment = false
)

// IsTestEnvironment should indicate whether the current run is a test run.
func IsTestEnvironment() bool {
	return isTestEnvironment
}

// SetIsTestEnvironment sets a flag that indicates that we are currently in
// a test environment.
func SetIsTestEnvironment(isTestEnv bool) {
	isTestEnvironment = isTestEnv
}

// Assert will throw a panic if condition is false.
// The additional parameter is provided to panic() as argument.
func Assert(cond bool, i interface{}) {
	if cond == false {
		panic(i)
	}
}

// AssertNoError is a cheap function to check for err==nil and
// panics if condition is not met
func AssertNoError(err error) {
	Assert(err == nil, err)
}

// ExitOnError exits the program with rc=1 if err!=nil. The following
// Message and all additional arguments will be passed to fmt.Printf()
func ExitOnError(err error, errMsg string, v ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg+": ", v...)
		fmt.Fprintln(os.Stderr, err)
		// TODO add flag somewhere that will exit with panic() instead
		os.Exit(1)
	}
}

// ExitWithMessage prints out a message on StdErr and then exits the
// program with rc=1
func ExitWithMessage(errMsg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+errMsg, v...)
	os.Exit(1)
}

// InformOnError provides an error message on stdErr if an error
// occurs, but continues with the program
func InformOnError(err error, errMsg string, v ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, errMsg+":\n", v...)
		fmt.Fprintln(os.Stderr, err)
	}
}

// GetTempDir returns the location of a temporary directory in which
// files can be stored. The caller needs to ensure that the
// directory is deleted afterwards.
func GetTempDir() (dirName string) {
	dirName, err := ioutil.TempDir("", "pfscf-")
	AssertNoError(err)
	return dirName
}

// GetExecutableDir returns the dir in which the binary of this program is located
func GetExecutableDir() (dirName string) {
	var baseDir string

	if IsTestEnvironment() {
		// during test runs the executable will be run from some temporary
		// directory, so instead return the local directory for that case.
		baseDir = "."
	} else {
		baseDir = filepath.Dir(os.Args[0])
	}

	dirName, err := filepath.Abs(baseDir)
	AssertNoError(err)
	return dirName
}

// IsFile checks wether a file exists and is not a directory.
func IsFile(filename string) (exists bool, err error) {
	info, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("Path %v is a directory", filename)
	}
	return true, nil
}

// IsDir checks wether a directory exists.
func IsDir(dirname string) (exists bool, err error) {
	info, err := os.Stat(dirname)
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("Path %v is not a directory", dirname)
	}
	return true, nil
}

// IsSet checks whether the provided value is different from its zero value and
// in case of a non-nil pointer it also checks whether the referenced value is
// not the zero value
func IsSet(val interface{}) (result bool) {
	if reflect.TypeOf(val) == nil {
		return false
	}
	x := reflect.ValueOf(val)
	if x.Kind() == reflect.Ptr {
		return !(x.IsNil() || x.Elem().IsZero())
	}
	return !x.IsZero()
}

// HasWhitespace checks if a string includes at least one whitespace character
func HasWhitespace(s string) bool {
	return strings.Index(s, " ") > -1
}

// QuoteStringIfRequired takes a string as input. If the string contains whitespace,
// then it is returned enclosed by quotation marks, else it is returned as-is
func QuoteStringIfRequired(s string) string {
	if HasWhitespace(s) {
		return fmt.Sprintf("\"%v\"", s)
	}
	return s
}

// IsExported takes a reflection value and checks whether that is an exported field in
// a struct. Well, we cannot check here whether this is a struct field, but under the
// assumption that it was one, we can check whether it is exported then. Yay!
func IsExported(val reflect.Value) bool {
	// Check to CanInterface() should be sufficient according to
	// https://stackoverflow.com/questions/50279840/when-is-go-reflect-caninterface-false
	return val.CanInterface()
}

// Contains checks whether a list of strings contains a specific string.
func Contains(list []string, element string) (result bool) {
	for _, listElement := range list {
		if element == listElement {
			return true
		}
	}
	return false
}

// SortCoords takes two coordinates and returns them ordered from large to small
func SortCoords(c1, c2 float64) (float64, float64) {
	if c1 > c2 {
		return c1, c2
	}
	return c2, c1
}

// AddMissingValues iterates over the exported fields of the source object. For each
// such fields it checks whether the target object contains a field with the same
// name. If that is the case and if the target field does not yet have a value set,
// then the value from the source object is copied over.
func AddMissingValues(target interface{}, source interface{}, ignoredFields ...string) {
	Assert(reflect.ValueOf(source).Kind() == reflect.Struct, "Can only process structs as source")
	Assert(reflect.ValueOf(target).Kind() == reflect.Ptr, "Target argument must be passed by ptr, as we modify it")
	Assert(reflect.ValueOf(target).Elem().Kind() == reflect.Struct, "Can only process structs as target")

	vSrc := reflect.ValueOf(source)
	vDst := reflect.ValueOf(target).Elem()

	for i := 0; i < vDst.NumField(); i++ {
		fieldDst := vDst.Field(i)
		fieldName := vDst.Type().Field(i).Name

		// Ignore the Presets field, as we do not want to take over values for this.
		if Contains(ignoredFields, fieldName) { // especially filter out "Presets" and "ID"
			continue
		}

		// take care to skip unexported fields
		if !fieldDst.CanSet() {
			continue
		}

		fieldSrc := vSrc.FieldByName(fieldName)

		// skip target fields that do not exist on source side side
		if !fieldSrc.IsValid() {
			continue
		}

		if fieldDst.IsZero() && !fieldSrc.IsZero() {
			switch fieldDst.Kind() {
			case reflect.String:
				fallthrough
			case reflect.Float64:
				fieldDst.Set(fieldSrc)
			default:
				panic(fmt.Sprintf("Unsupported datat type '%v' in struct, update function 'AddMissingValuesFrom()'", fieldDst.Kind()))
			}
		}
	}
}

// ReadFileToLines reads in the given file and returns the content as string array of lines
func ReadFileToLines(filename string) (lines []string, err error) {
	lines = make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineScanner := bufio.NewScanner(file)
	for lineScanner.Scan() {
		lines = append(lines, lineScanner.Text())
	}

	if lineScanner.Err() != nil {
		return nil, lineScanner.Err()
	}

	return lines, nil
}

// OpenWithDefaultViewer opens the provided file with the default viewer registered in the system for this.
func OpenWithDefaultViewer(file string) (err error) {
	absFile, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", absFile).Start()
	case "windows":
		err = exec.Command(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe"), "url.dll,FileProtocolHandler", absFile).Start()
	case "darwin":
		err = exec.Command("open", absFile).Start()
	default:
		err = fmt.Errorf("Unknown OS, cannot open file '%v' automatically", file)
	}

	return err
}
