package encode

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// ConvertByteToUtf8 checks whether the provided input is UTF8 encoded. If that is
// not the case, it will assume that the input is instead encoded in ISO 8859-1 / cp1252
// and convert this to UTF8
func ConvertByteToUtf8(input []byte) (output []byte, err error) {
	if utf8.Valid(input) {
		return input, nil
	}

	output, err = charmap.ISO8859_1.NewDecoder().Bytes(input)
	return output, err
}

// ConvertStringToUtf8 checks whether the provided input is UTF8 encoded. If that is
// not the case, it will assume that the input is instead encoded in ISO 8859-1 / cp1252
// and convert this to UTF8
func ConvertStringToUtf8(input string) (output string, err error) {
	outputBytes, err := ConvertByteToUtf8([]byte(input))
	if err != nil {
		return "", err
	}
	return string(outputBytes), nil
}
