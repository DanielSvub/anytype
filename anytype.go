/*
AnyType Library for Go
Type enum, field interface, atomic types
*/

package anytype

import (
	"math"
	"strconv"
)

/*
Type for AnyType data types.
*/
type Type uint8

/*
Enum of AnyType data types.
*/
const (
	TypeUndefined Type = iota
	TypeNil
	TypeObject
	TypeList
	TypeString
	TypeBool
	TypeInt
	TypeFloat
)

/*
Interface representing an AnyType field. Allows nesting of lists and objects.
*/
type field interface {
	getVal() any
	copy() any
	serialize() string
	isEqual(another any) bool
}

/*
Converts a given value to an AnyType field.
Parameters:
  - val - value to convert

Returns:
  - field, ready to add to object or list
*/
func parseVal(val any) field {
	switch v := val.(type) {
	case Object:
		return v
	case map[string]any:
		return NewObjectFrom(v)
	case map[string]Object:
		return NewObjectFrom(v)
	case map[string]List:
		return NewObjectFrom(v)
	case map[string]string:
		return NewObjectFrom(v)
	case map[string]bool:
		return NewObjectFrom(v)
	case map[string]int:
		return NewObjectFrom(v)
	case map[string]float64:
		return NewObjectFrom(v)
	case List:
		return v
	case []any:
		return NewListFrom(v)
	case []Object:
		return NewListFrom(v)
	case []List:
		return NewListFrom(v)
	case []string:
		return NewListFrom(v)
	case []bool:
		return NewListFrom(v)
	case []int:
		return NewListFrom(v)
	case []float64:
		return NewListFrom(v)
	case string:
		return newString(v)
	case bool:
		return newBool(v)
	case int:
		return newInt(v)
	case int64:
		return newInt(int(v))
	case int32:
		return newInt(int(v))
	case int16:
		return newInt(int(v))
	case int8:
		return newInt(int(v))
	case uint64:
		return newInt(int(v))
	case uint32:
		return newInt(int(v))
	case uint16:
		return newInt(int(v))
	case uint8:
		return newInt(int(v))
	case float64:
		return newFloat(v)
	case float32:
		return newFloat(float64(v))
	case nil:
		return newNil()
	default:
		panic("Unknown type.")
	}
}

/*
Structure encapsulating a string value.
Implements:
  - Fielder.
*/
type atString struct {
	val string
}

/*
Crates a new AnyType string.
Parameters:
  - val - value of the string.

Returns:
  - Pointer to the created string.
*/
func newString(val string) *atString {
	obj := atString{val: val}
	return &obj
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case a string value (string is a value type).
Returns:
  - value of the field.
*/
func (ego *atString) getVal() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new string.
Returns:
  - deep copy of the field.
*/
func (ego *atString) copy() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case simply prints the value.
Returns:
  - string representing serialized field.
*/
func (ego *atString) serialize() string {
	val := ego.getVal().(string)
	return strconv.Quote(val)
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *atString) isEqual(another any) bool {
	str, ok := another.(*atString)
	if !ok {
		return false
	}
	return ego.val == str.val
}

/*
Structure encapsulating a boolean value.
Implements:
  - Fielder.
*/
type atBool struct {
	val bool
}

/*
Crates a new AnyType bool.
Parameters:
  - val - value of the bool.

Returns:
  - Pointer to the created bool.
*/
func newBool(val bool) *atBool {
	obj := atBool{val: val}
	return &obj
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case a bool value (bool is a value type).
Returns:
  - value of the field.
*/
func (ego *atBool) getVal() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints a string representation of the value.
Returns:
  - string representing serialized field.
*/
func (ego *atBool) serialize() string {
	return strconv.FormatBool(ego.getVal().(bool))
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new bool.
Returns:
  - deep copy of the field.
*/
func (ego *atBool) copy() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *atBool) isEqual(another any) bool {
	boolean, ok := another.(*atBool)
	if !ok {
		return false
	}
	return ego.val == boolean.val
}

/*
Structure encapsulating an integer value.
Implements:
  - Fielder.
*/
type atInt struct {
	val int
}

/*
Crates a new AnyType int.
Parameters:
  - val - value of the int.

Returns:
  - Pointer to the created int.
*/
func newInt(val int) *atInt {
	obj := atInt{val: val}
	return &obj
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case an int value (int is a value type).
Returns:
  - value of the field.
*/
func (ego *atInt) getVal() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints a string representation of the value.
Returns:
  - string representing serialized field.
*/
func (ego *atInt) serialize() string {
	return strconv.Itoa(ego.getVal().(int))
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new int.
Returns:
  - deep copy of the field.
*/
func (ego *atInt) copy() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *atInt) isEqual(another any) bool {
	integer, ok := another.(*atInt)
	if !ok {
		return false
	}
	return ego.val == integer.val
}

/*
Structure encapsulating a float value.
Implements:
  - Fielder.
*/
type atFloat struct {
	val float64
}

/*
Crates a new AnyType float.
Parameters:
  - val - value of the float.

Returns:
  - Pointer to the created float.
*/
func newFloat(val float64) *atFloat {
	obj := atFloat{val: val}
	return &obj
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case a float value (float is a value type).
Returns:
  - value of the field.
*/
func (ego *atFloat) getVal() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints a string representation of the value.
Returns:
  - string representing serialized field.
*/
func (ego *atFloat) serialize() string {
	val := ego.getVal().(float64)
	abs := math.Abs(val)
	if abs >= math.Pow10(6) || (abs > 0 && abs <= math.Pow10(-6)) {
		return strconv.FormatFloat(val, 'e', -1, 64)
	}
	return strconv.FormatFloat(val, 'f', -1, 64)
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new float.
Returns:
  - deep copy of the field.
*/
func (ego *atFloat) copy() any {
	return ego.val
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *atFloat) isEqual(another any) bool {
	float, ok := another.(*atFloat)
	if !ok {
		return false
	}
	return ego.val == float.val
}

/*
Structure encapsulating a nil value.
Implements:
  - Fielder.
*/
type atNil struct {
}

/*
Crates a new AnyType nil.
Returns:
  - Pointer to the created nil.
*/
func newNil() *atNil {
	return &atNil{}
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case nil.
Returns:
  - value of the field.
*/
func (ego *atNil) getVal() any {
	return nil
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints "null".
Returns:
  - string representing serialized field.
*/
func (ego *atNil) serialize() string {
	return "null"
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new nil.
Returns:
  - deep copy of the field.
*/
func (ego *atNil) copy() any {
	return nil
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *atNil) isEqual(another any) bool {
	_, ok := another.(*atNil)
	return ok
}
