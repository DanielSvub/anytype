/*
AnyType Library for Go
JSON parser
*/

package anytype

import (
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
const (
	stateStart = iota
	stateKeyStart
	stateKey
	stateKeyEscape
	stateAfterKey
	stateVal
	stateValEscape
	stateValString
	stateValAfterString
)

/*
Parses a primitive field of object or list. Does not include strings.
Parameters:
  - field - field to parse.

Returns:
  - parsed value (nil, int, float64 or bool).
*/
func parseField(field string) any {
	if field == "null" {
		return nil
	}
	integer, err := strconv.ParseInt(field, 0, bits.UintSize)
	if err == nil {
		return int(integer)
	}
	float, err := strconv.ParseFloat(field, bits.UintSize)
	if err == nil {
		return float
	}
	boolean, err := strconv.ParseBool(field)
	if err == nil {
		return boolean
	}
	panic("Not a valid JSON - Invalid value '" + field + "'.")
}

/*
Recursively parses a JSON list.
Patameters:
  - json - JSON string to parse.

Returns:
  - created list,
  - number of chars processed.
*/
func parseList(json string) (List, int) {

	state := stateStart
	var list List
	var val string
	var inVal bool

	var char rune
	var size int

	// Iterating over all characters
	for i := 0; i < len(json); i += size {

		char, size = utf8.DecodeLastRuneInString(json)
		if size == 0 {
			panic("Invalid encoding.")
		}

		switch state {

		// List creation
		case stateStart:
			if char != '[' {
				panic("Not a valid JSON - Expecting '{', got '" + string(char) + "'.")
			}
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
				o, pos := parseObject(json[i:])
				i += pos
				list.Add(o)
				continue
			}

			// Nested list (same as above)
			if !inVal && char == '[' {
				l, pos := parseList(json[i:])
				i += pos
				list.Add(l)
				continue
			}

			// End of the element
			if char == ',' || char == ']' {
				if len(val) > 0 {
					list.Add(parseField(val))
					val = ""
					inVal = false
				}
				if char == ']' {
					return list, i
				}
				continue
			}

			// Inside the element
			val += string(char)
			inVal = true

		// Parsing a string
		case stateValString:
			if char == '\\' {
				state = stateValEscape
				continue
			}
			if char == '"' {
				list.Add(val)
				val = ""
				state = stateValAfterString
				continue
			}
			val += string(char)

		// Escaping inside string
		case stateValEscape:
			val += string(char)
			state = stateValString

		// End of the string
		case stateValAfterString:
			if char == ',' {
				state = stateVal
				continue
			} else if char == ']' {
				return list, i
			}

		}
	}

	// No matching rule - panic
	panic("Not a valid JSON - Unexpected end of input.")

}

/*
Recursively parses a JSON object.
Patameters:
  - json - JSON string to parse.

Returns:
  - created object,
  - number of chars processed.
*/
func parseObject(json string) (Object, int) {

	state := stateStart
	var object Object
	var key string
	var val string
	var inVal bool

	// Iterating over all characters
	for i := 0; i < len(json); i++ {

		char := rune(json[i])

		switch state {

		// List creation
		case stateStart:
			if char != '{' {
				panic("Not a valid JSON - Expecting '{', got '" + string(char) + "'.")
			}
			object = NewObject()
			state = stateKeyStart

		// Begining of a key
		case stateKeyStart:
			if unicode.IsSpace(char) {
				continue
			}
			if char == '}' {
				return object, i
			}
			if char == '"' {
				key = ""
				state = stateKey
				continue
			}
			panic(`Not a valid JSON - Expecting '"', got '` + string(char) + `'.`)

		// Parsing the key
		case stateKey:
			if char == '"' {
				state = stateAfterKey
			} else if char == '\\' {
				state = stateKeyEscape
			} else {
				key += string(char)
			}

		// Escaping inside key
		case stateKeyEscape:
			key += string(char)
			state = stateKey

		// Waiting for colon
		case stateAfterKey:
			if unicode.IsSpace(char) {
				continue
			}
			if char != ':' {
				panic("Not a valid JSON - Expecting ':', got '" + string(char) + "'.")
			}
			val = ""
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
				o, pos := parseObject(json[i:])
				i += pos
				object.Set(key, o)
				continue
			}

			// Nested list (same as above)
			if !inVal && char == '[' {
				l, pos := parseList(json[i:])
				i += pos
				object.Set(key, l)
				continue
			}

			// End of the value
			if char == ',' || char == '}' {
				if len(val) > 0 {
					object.Set(key, parseField(val))
				}
				if char == ',' {
					state = stateKeyStart
					continue
				} else if char == '}' {
					return object, i
				}
			}

			// Inside the value
			val += string(char)
			inVal = true

		// Parsing a string
		case stateValString:
			if char == '\\' {
				state = stateValEscape
				continue
			}
			if char == '"' {
				object.Set(key, val)
				state = stateValAfterString
				continue
			}
			val += string(char)

		// Escaping inside string
		case stateValEscape:
			val += string(char)
			state = stateValString

		// End of the string
		case stateValAfterString:
			if char == ',' {
				state = stateKeyStart
				continue
			} else if char == '}' {
				return object, i
			}

		}
	}

	// No matching rule - panic
	panic("Not a valid JSON - Unexpected end of input.")

}

/*
Creates a new object from JSON.
Patameters:
  - json - JSON string to parse.

Returns:
  - created object.
*/
func ParseJson(json string) Object {
	root, _ := parseObject(strings.TrimSpace(json))
	return root
}

/*
Creates a new object from JSON file.
Patameters:
  - path - path to the file to parse.

Returns:
  - created object.
*/
func ParseFile(path string) Object {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return ParseJson(string(data))
}
