/*
AnyType Library for Go
JSON parser
*/

package anytype

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
Enum of possible states of the parser.
*/
type parserState uint8

const (
	stateStart parserState = iota
	stateKeyStart
	stateKey
	stateKeyEscape
	stateAfterKey
	stateVal
	stateAfterVal
	stateValEscape
	stateValString
	stateValAfterString
)

/*
Parses a primitive field of object or list. Does not include strings.
Parameters:
  - field - field to parse.

Returns:
  - parsed value (nil, int, float64 or bool),
  - error if any occurred.
*/
func parseField(field string, line int) (any, error) {
	if field == "null" {
		return nil, nil
	}
	integer, err := strconv.ParseInt(field, 0, bits.UintSize)
	if err == nil {
		return int(integer), nil
	}
	float, err := strconv.ParseFloat(field, bits.UintSize)
	if err == nil {
		return float, nil
	}
	boolean, err := strconv.ParseBool(field)
	if err == nil {
		return boolean, nil
	}
	return nil, fmt.Errorf("not a valid JSON - invalid value '%s' on line %d", field, line)
}

/*
Recursively parses a JSON list.
Patameters:
  - json - JSON string to parse,
  - line - current line of the input.

Returns:
  - created list,
  - number of bytes processed,
  - error if any occurred.
*/
func parseList(json string, line *int) (List, int, error) {

	state := stateStart
	var list List
	var val strings.Builder
	var inVal bool

	var char rune
	var size int

	// Iterating over all characters
	for i := 0; i < len(json); i += size {

		char, size = utf8.DecodeRuneInString(json[i:])
		if size == 0 || char == utf8.RuneError {
			return nil, 0, fmt.Errorf("not an UTF-8 encoding")
		}

		if char == '\n' {
			*line++
		}

		switch state {

		// List creation
		case stateStart:
			list = NewList()
			state = stateVal
			inVal = false

		// Parsing an element
		case stateVal:

			// Whitespace (skipping)
			if unicode.IsSpace(char) {
				continue
			}

			// Start of a string
			if !inVal && char == '"' {
				state = stateValString
				continue
			}

			// Nested object
			// Recursive call with original string starting from current position
			// Current index is moved after the nested object so the parsing can continue
			if !inVal && char == '{' {
				o, pos, err := parseObject(json[i:], line)
				if err != nil {
					return nil, 0, err
				}
				i += pos
				list.Add(o)
				continue
			}

			// Nested list (same as above)
			if !inVal && char == '[' {
				l, pos, err := parseList(json[i:], line)
				if err != nil {
					return nil, 0, err
				}
				i += pos
				list.Add(l)
				continue
			}

			// End of the element
			if char == ',' || char == ']' {
				if val.Len() > 0 {
					field, err := parseField(val.String(), *line)
					if err != nil {
						return nil, 0, err
					}
					list.Add(field)
					val.Reset()
					inVal = false
				}
				if char == ']' {
					return list, i, nil
				}
				continue
			}

			// Inside the element
			val.WriteRune(char)
			inVal = true

		// Parsing a string
		case stateValString:
			if char == '\\' {
				state = stateValEscape
				continue
			}
			if char == '"' {
				str, _ := strconv.Unquote(fmt.Sprintf(`"%s"`, val.String()))
				list.Add(str)
				val.Reset()
				state = stateValAfterString
				continue
			}
			val.WriteRune(char)

		// Escaping inside string
		case stateValEscape:
			val.WriteRune('\\')
			val.WriteRune(char)
			state = stateValString

		// End of the string
		case stateValAfterString:
			if char == ',' {
				state = stateVal
				continue
			} else if char == ']' {
				return list, i, nil
			}

		}
	}

	// No matching rule - error
	return nil, 0, fmt.Errorf("not a valid JSON - unexpected end of input")

}

/*
Recursively parses a JSON object.
Patameters:
  - json - JSON string to parse,
  - line - current line of the input.

Returns:
  - created object,
  - number of bytes processed,
  - error if any occurred.
*/
func parseObject(json string, line *int) (Object, int, error) {

	state := stateStart
	var object Object
	var key strings.Builder
	var val strings.Builder
	var inVal bool

	var char rune
	var size int

	// Iterating over all characters
	for i := 0; i < len(json); i += size {

		char, size = utf8.DecodeRuneInString(json[i:])
		if size == 0 || char == utf8.RuneError {
			return nil, 0, fmt.Errorf("not an UTF-8 encoding")
		}

		if char == '\n' {
			*line++
		}

		switch state {

		// Object creation
		case stateStart:
			object = NewObject()
			state = stateKeyStart

		// Begining of a key
		case stateKeyStart:
			if unicode.IsSpace(char) {
				continue
			}
			if char == '}' {
				return object, i, nil
			}
			if char == '"' {
				key.Reset()
				state = stateKey
				continue
			}
			return nil, 0, fmt.Errorf("not a valid JSON - expecting '\"', got '%s' on line %d", string(char), *line)

		// Parsing the key
		case stateKey:
			if char == '"' {
				state = stateAfterKey
			} else if char == '\\' {
				state = stateKeyEscape
			} else {
				key.WriteRune(char)
			}

		// Escaping inside key
		case stateKeyEscape:
			key.WriteRune('\\')
			key.WriteRune(char)
			state = stateKey

		// Waiting for colon
		case stateAfterKey:
			if unicode.IsSpace(char) {
				continue
			}
			if char != ':' {
				return nil, 0, fmt.Errorf("not a valid JSON - expecting ':', got '%s' on line %d", string(char), *line)
			}
			str, _ := strconv.Unquote(fmt.Sprintf(`"%s"`, key.String()))
			key.Reset()
			key.WriteString(str)
			val.Reset()
			state = stateVal
			inVal = false

		// Parsing a value
		case stateVal:

			// Whitespace (skipping)
			if unicode.IsSpace(char) {
				continue
			}

			// Start of a string
			if !inVal && char == '"' {
				state = stateValString
				continue
			}

			// Nested object
			// Recursive call with original string starting from current position
			// Current index is moved after the nested object so the parsing can continue
			if !inVal && char == '{' {
				o, pos, err := parseObject(json[i:], line)
				if err != nil {
					return nil, 0, err
				}
				i += pos
				object.Set(key.String(), o)
				state = stateAfterVal
				continue
			}

			// Nested list (same as above)
			if !inVal && char == '[' {
				l, pos, err := parseList(json[i:], line)
				if err != nil {
					return nil, 0, err
				}
				i += pos
				object.Set(key.String(), l)
				state = stateAfterVal
				continue
			}

			// End of the value
			if char == ',' || char == '}' {
				if val.Len() > 0 {
					field, err := parseField(val.String(), *line)
					if err != nil {
						return nil, 0, err
					}
					object.Set(key.String(), field)
				}
				if char == ',' {
					state = stateKeyStart
					continue
				} else if char == '}' {
					return object, i, nil
				}
			}

			// Inside the value
			val.WriteRune(char)
			inVal = true

		// After nested object or list
		case stateAfterVal:

			// Whitespace (skipping)
			if unicode.IsSpace(char) {
				continue
			}

			if char == ',' || char == '}' {
				if char == ',' {
					state = stateKeyStart
					continue
				}
				return object, i, nil
			}

			if char == '"' {
				key.Reset()
				state = stateKey
				continue
			}

			return nil, 0, fmt.Errorf("not a valid JSON - expecting ',' or '}', got '%s' on line %d", string(char), *line)

		// Parsing a string
		case stateValString:
			if char == '\\' {
				state = stateValEscape
				continue
			}
			if char == '"' {
				str, _ := strconv.Unquote(fmt.Sprintf(`"%s"`, val.String()))
				object.Set(key.String(), str)
				state = stateValAfterString
				continue
			}
			val.WriteRune(char)

		// Escaping inside string
		case stateValEscape:
			val.WriteRune('\\')
			val.WriteRune(char)
			state = stateValString

		// End of the string
		case stateValAfterString:
			if char == ',' {
				state = stateKeyStart
				continue
			} else if char == '}' {
				return object, i, nil
			}

		}
	}

	// No matching rule - error
	return nil, 0, fmt.Errorf("not a valid JSON - unexpected end of input")

}

/*
ParseList creates a new list from JSON.
Patameters:
  - json - JSON string to parse.

Returns:
  - created list,
  - error if any occurred.
*/
func ParseList(json string) (List, error) {
	start := strings.Index(json, "[")
	if start < 0 {
		return nil, fmt.Errorf("not a valid JSON - missing '['")
	}
	startLine := strings.Count(json[:start], "\n") + 1
	root, _, err := parseList(json[start:], &startLine)
	return root, err
}

/*
ParseObject creates a new object from JSON.
Patameters:
  - json - JSON string to parse.

Returns:
  - created object,
  - error if any occurred.
*/
func ParseObject(json string) (Object, error) {
	start := strings.Index(json, "{")
	if start < 0 {
		return nil, fmt.Errorf("not a valid JSON - missing '{'")
	}
	startLine := strings.Count(json[:start], "\n") + 1
	root, _, err := parseObject(json[start:], &startLine)
	return root, err
}

/*
ParseFile creates a new object from JSON file.
Patameters:
  - path - path to the file to parse.

Returns:
  - created object,
  - error if any occurred.
*/
func ParseFile(path string) (Object, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseObject(string(data))
}
