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
	fmt.Fprintf(os.Stderr, "Error: "+errMsg+"\n", v...)
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

// IsSet checks whether the provided value is different from its zero value for
// element types and different from nil for pointer type.s
func IsSet(val interface{}) (result bool) {
	if reflect.TypeOf(val) == nil {
		return false
	}
	x := reflect.ValueOf(val)
	if x.Kind() == reflect.Ptr {
		return !x.IsNil()
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

// UnquoteStringIfRequired takes a string as input. If the string is enclosed in
// quotation marks, then these are removed in the returned string
func UnquoteStringIfRequired(s string) string {
	if len(s) > 1 {
		first := s[:1]
		last := s[len(s)-1:]
		if first == last && Contains([]string{"\"", "'"}, first) {
			return s[1 : len(s)-1]
		}
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

func copyValueIfUnset(vSrc, vDst reflect.Value) {
	Assert(vDst.IsValid(), "Valid destination value expected")
	Assert(vDst.CanSet(), "Settable destination value expected")
	Assert(vSrc.IsValid(), "Valid source value expected")

	introspect := func(v reflect.Value) (isSet, isPtr bool) {
		if v.Kind() == reflect.Ptr {
			return !v.IsNil(), true
		}
		return !v.IsZero(), false
	}

	srcIsSet, srcIsPtr := introspect(vSrc)
	dstIsSet, dstIsPtr := introspect(vDst)

	if !dstIsSet && srcIsSet {
		var dstElem reflect.Value
		if dstIsPtr {
			// dst is a nil pointer, so allocate space and assign

			// get required elem type from source
			var srcElemType reflect.Type
			if srcIsPtr {
				srcElemType = vSrc.Elem().Type()
			} else {
				srcElemType = vSrc.Type()
			}

			dstPtr := reflect.New(srcElemType)
			vDst.Set(dstPtr)
			dstElem = vDst.Elem()
		} else {
			dstElem = vDst
		}

		var srcElem reflect.Value
		if srcIsPtr {
			srcElem = vSrc.Elem()
		} else {
			srcElem = vSrc
		}

		switch dstElem.Kind() {
		case reflect.String:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Float64:
			dstElem.Set(srcElem)
		default:
			panic(fmt.Sprintf("Unsupported data type '%v' in struct, update function 'AddMissingValuesFrom()'", dstElem.Kind()))
		}
	}

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

		// Ignore some fields, as we do not want to take over values for this.
		if Contains(ignoredFields, fieldName) { // especially filter out "Presets" and "ID"
			continue
		}

		// take care to skip unexported fields
		if !fieldDst.CanSet() {
			continue
		}

		fieldSrc := vSrc.FieldByName(fieldName)

		// skip target fields that do not exist on source side
		if !fieldSrc.IsValid() {
			continue
		}

		copyValueIfUnset(fieldSrc, fieldDst)

		fctIsSet := func(v reflect.Value) bool {
			if v.Kind() == reflect.Ptr {
				return !v.IsNil()
			}
			return !v.IsZero()
		}

		Assert(!fctIsSet(vDst) || fctIsSet(vSrc), "This did not work as expected")
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

// GenericFieldsCheck provides a generic way to check multiple fields of a struct at once.
func GenericFieldsCheck(obj interface{}, isOk func(interface{}) bool, fieldNames ...string) (err error) {
	oVal := reflect.ValueOf(obj)
	if oVal.Kind() == reflect.Ptr {
		oVal = oVal.Elem()
	}
	Assert(oVal.Kind() == reflect.Struct, "Can only work on structs or pointers to structs")

	errFields := make([]string, 0)
	for _, fieldName := range fieldNames {
		fieldVal := oVal.FieldByName(fieldName)
		Assert(fieldVal.IsValid(), fmt.Sprintf("No field with name '%v' found in struct of type '%T'", fieldName, obj))

		if !isOk(fieldVal.Interface()) {
			errFields = append(errFields, fieldName)
		}
	}

	if len(errFields) > 0 {
		return fmt.Errorf("%v", errFields)
	}
	return nil
}

// CheckFieldsAreSet checks whether all provided struct fields have a non-zero value.
func CheckFieldsAreSet(obj interface{}, fieldNames ...string) (err error) {
	err = GenericFieldsCheck(obj, IsSet, fieldNames...)
	if err != nil {
		return fmt.Errorf("Missing values for the following fields: %v", err)
	}
	return nil
}

// CheckFieldsAreInRange checks that all provided struct fields are within a given range.
func CheckFieldsAreInRange(obj interface{}, min, max float64, fieldNames ...string) (err error) {
	isOk := func(obj interface{}) bool {
		oVal := reflect.ValueOf(obj)
		if oVal.Kind() == reflect.Ptr {
			oVal = oVal.Elem()
		}
		fObj := oVal.Interface().(float64)
		return fObj >= min && fObj <= max
	}
	err = GenericFieldsCheck(obj, isOk, fieldNames...)
	if err != nil {
		return fmt.Errorf("Values for the following fields are out of range %.2f-%.2f: %v", min, max, err)
	}
	return nil
}

// CopyString creates a copy of a referenced string and returns the result as pointer.
func CopyString(in *string) (out *string) {
	if in == nil {
		return nil
	}
	out = new(string)
	*out = *in
	return out
}

// CopyFloat creates a copy of a referenced float and returns the result as pointer.
func CopyFloat(in *float64) (out *float64) {
	if in == nil {
		return nil
	}
	out = new(float64)
	*out = *in
	return out
}

// ToCommaSeparatedString takes a list of input values and returns them as single
// string with all values separated by commas.
func ToCommaSeparatedString(input []string) (result string) {
	var sb strings.Builder

	for index, current := range input {
		fmt.Fprint(&sb, current)
		if index != len(input)-1 {
			fmt.Fprint(&sb, ", ")
		}
	}

	return sb.String()
}

// SplitAndTrim splits the string up at the provided separator and them trims each substring.
func SplitAndTrim(input, sep string) (result []string) {
	result = strings.Split(input, sep)

	for idx, entry := range result {
		result[idx] = strings.TrimSpace(entry)
	}

	return result
}
