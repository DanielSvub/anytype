/*
AnyType Library for Go
Object (dictionary) type
*/

package anytype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

/*
Interface for an object.

Extends:
  - field.
*/
type Object interface {
	field

	Init(ptr Object)

	// Manipulation with fields
	Set(values ...any) Object
	Unset(keys ...string) Object
	Clear() Object

	// Getting fields
	Get(key string) any
	GetObject(key string) Object
	GetList(key string) List
	GetString(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool

	// TypeOf check
	TypeOf(key string) Type

	// Export
	String() string
	FormatString(indent int) string
	Dict() map[string]any
	Keys() List
	Values() List

	// Features over whole object
	Clone() Object
	Count() int
	Empty() bool
	Equals(another Object) bool
	Merge(another Object) Object
	Pluck(keys ...string) Object
	Contains(elem any) bool
	KeyOf(elem any) string
	KeyExists(key string) bool

	// ForEaches
	ForEach(function func(string, any)) Object
	ForEachValue(function func(any)) Object
	ForEachObject(function func(Object)) Object
	ForEachList(function func(List)) Object
	ForEachString(function func(string)) Object
	ForEachBool(function func(bool)) Object
	ForEachInt(function func(int)) Object
	ForEachFloat(function func(float64)) Object

	// Mapping
	Map(function func(string, any) any) Object
	MapValues(function func(any) any) Object
	MapObjects(function func(Object) any) Object
	MapLists(function func(List) any) Object
	MapStrings(function func(string) any) Object
	MapInts(function func(int) any) Object
	MapFloats(function func(float64) any) Object

	// Async
	ForEachAsync(function func(string, any)) Object
	MapAsync(function func(string, any) any) Object

	// Tree form
	GetTF(tf string) any
	SetTF(tf string, value any) Object
}

/*
Map object, a reference type. Contains a map.

Implements:
  - Fielder,
  - Objecter.
*/
type MapObject struct {
	val map[string]field
	ptr Object
}

/*
Object constructor.
Creates a new object.

Parameters:
  - values... - any amount of key-value pairs to set after the object creation.

Returns:
  - pointer to the created object.
*/
func NewObject(vals ...any) Object {
	ego := &MapObject{val: make(map[string]field)}
	ego.ptr = ego
	ego.Set(vals...)
	return ego
}

/*
Object constructor.
Converts a map of supported types to an object.

Parameters:
  - dict - original map.

Returns:
  - created object.
*/
func NewObjectFrom(dict any) Object {
	object := NewObject()
	switch s := dict.(type) {
	case map[string]any:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]Object:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]List:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]string:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]bool:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]int:
		for key, value := range s {
			object.Set(key, value)
		}
	case map[string]float64:
		for key, value := range s {
			object.Set(key, value)
		}
	default:
		panic("Unknown map type.")
	}
	return object
}

/*
Asserts that the object is initialized.
*/
func (ego *MapObject) assert() {
	if ego == nil || ego.val == nil {
		panic("Object is not initialized.")
	}
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case a reference to the whole struct (Object is reference type).

Returns:
  - value of the field.
*/
func (ego *MapObject) getVal() any {
	return ego.ptr
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new object with identical fields.
Can be called recursively.

Returns:
  - deep copy of the field.
*/
func (ego *MapObject) copy() any {
	obj := NewObject()
	for key, value := range ego.val {
		obj.Set(key, value.copy())
	}
	return obj
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints all keys and their values.
Can be called recursively.

Returns:
  - string representing serialized field.
*/
func (ego *MapObject) serialize() string {
	result := "{"
	i := 0
	for field, value := range ego.val {
		result += fmt.Sprintf("%s:%s", strconv.Quote(field), value.serialize())
		if i++; i < len(ego.val) {
			result += ","
		}
	}
	result += "}"
	return result
}

/*
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Can be called recursively.

Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *MapObject) isEqual(another any) bool {
	obj, ok := another.(*MapObject)
	if !ok || ego.Count() != obj.Count() {
		return false
	}
	for k := range ego.val {
		if !ego.val[k].isEqual(obj.val[k]) {
			return false
		}
	}
	return true
}

/*
Initializes the ego pointer, which allows deriving.

Parameters:
  - ptr - ego pointer.
*/
func (ego *MapObject) Init(ptr Object) {
	ego.ptr = ptr
}

/*
Sets a values of the fields.
If the key already exists, the value is overwritten, if not, new field is created.
If one key is given multiple times, the value is set to the last one.

Parameters:
  - values... - any amount of key-value pairs to set.

Returns:
  - updated object.
*/
func (ego *MapObject) Set(values ...any) Object {
	ego.assert()
	length := len(values)
	if length%2 != 0 {
		panic("Object fields have to be set as key-value pairs.")
	}
	for i := 0; i < length; i += 2 {
		name, ok := values[i].(string)
		if !ok || name == "" {
			panic("Object key has to be non-empty string.")
		}
		ego.val[name] = parseVal(values[i+1])
	}
	return ego.ptr
}

/*
Deletes the fields with given keys.

Parameters:
  - keys... - any amount of keys to delete.

Returns:
  - updated object.
*/
func (ego *MapObject) Unset(keys ...string) Object {
	ego.assert()
	for _, key := range keys {
		if ego.val[key] == nil {
			panic("Object does not have a field '" + key + "'.")
		}
		delete(ego.val, key)
	}
	return ego.ptr
}

/*
Deletes all field of the object.

Returns:
  - updated object.
*/
func (ego *MapObject) Clear() Object {
	ego.assert()
	ego.val = make(map[string]field, 0)
	return ego.ptr
}

/*
Acquires the value under the specified key of the object.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value (any type, has to be asserted).
*/
func (ego *MapObject) Get(key string) any {
	ego.assert()
	if ego.val[key] == nil {
		panic("Object does not have a field '" + key + "'.")
	}
	obj := ego.val[key]
	switch obj.(type) {
	case Object, List:
		return obj
	default:
		return obj.getVal()
	}
}

/*
Acquires the nested object under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as object.
*/
func (ego *MapObject) GetObject(key string) Object {
	o, ok := ego.Get(key).(*MapObject)
	if !ok {
		panic("Item is not an object.")
	}
	return o
}

/*
Acquires the list under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as list.
*/
func (ego *MapObject) GetList(key string) List {
	o, ok := ego.Get(key).(*SliceList)
	if !ok {
		panic("Field is not a list.")
	}
	return o
}

/*
Acquires the string under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as string.
*/
func (ego *MapObject) GetString(key string) string {
	o, ok := ego.Get(key).(string)
	if !ok {
		panic("Field is not a string.")
	}
	return o
}

/*
Acquires the bool under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as bool.
*/
func (ego *MapObject) GetBool(key string) bool {
	o, ok := ego.Get(key).(bool)
	if !ok {
		panic("Field is not a bool.")
	}
	return o
}

/*
Acquires the int under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as int.
*/
func (ego *MapObject) GetInt(key string) int {
	o, ok := ego.Get(key).(int)
	if !ok {
		panic("Field is not an int.")
	}
	return o
}

/*
Acquires the float under the specified key of the object.
Causes a panic if the field has another type.

Parameters:
  - key - key of the field to get.

Returns:
  - corresponding value asserted as float.
*/
func (ego *MapObject) GetFloat(key string) float64 {
	o, ok := ego.Get(key).(float64)
	if !ok {
		panic("Field is not a float.")
	}
	return o
}

/*
Gives a type of the field under the specified key of the object.

Parameters:
  - key - key of the field.

Returns:
  - integer constant representing the type (see type enum).
*/
func (ego *MapObject) TypeOf(key string) Type {
	ego.assert()
	switch ego.val[key].(type) {
	case *atString:
		return TypeString
	case *atInt:
		return TypeInt
	case *atBool:
		return TypeBool
	case *atFloat:
		return TypeFloat
	case *MapObject:
		return TypeObject
	case *SliceList:
		return TypeList
	case *atNil:
		return TypeNil
	default:
		panic("Unknown field type.")
	}
}

/*
Gives a JSON representation of the object, including nested objects and lists.

Returns:
  - JSON string.
*/
func (ego *MapObject) String() string {
	ego.assert()
	return ego.ptr.serialize()
}

/*
Gives a JSON representation of the object in standardized format with the given indentation.

Parameters:
  - indent - indentation spaces (0-10).

Returns:
  - JSON string.
*/
func (ego *MapObject) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic("Invalid indentation.")
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

/*
Converts the object into a Go map of empty interfaces.

Returns:
  - map.
*/
func (ego *MapObject) Dict() map[string]any {
	ego.assert()
	dict := make(map[string]any, 0)
	for key, value := range ego.val {
		dict[key] = value.getVal()
	}
	return dict
}

/*
Convers the object to a list of its keys.

Returns:
  - list of keys of the object.
*/
func (ego *MapObject) Keys() List {
	keys := NewList()
	for key := range ego.val {
		keys.Add(key)
	}
	return keys
}

/*
Convers the object to a list of its values.

Returns:
  - list of values of the object.
*/
func (ego *MapObject) Values() List {
	values := NewList()
	for _, value := range ego.val {
		values.Add(value.getVal())
	}
	return values
}

/*
Creates a deep copy of the object.

Returns:
  - copied object.
*/
func (ego *MapObject) Clone() Object {
	ego.assert()
	return ego.ptr.copy().(*MapObject)
}

/*
Gives a number of fields of the object.

Returns:
  - number of fields.
*/
func (ego *MapObject) Count() int {
	ego.assert()
	return len(ego.val)
}

/*
Checks whether the object is empty.

Returns:
  - true if the object is empty, false otherwise.
*/
func (ego *MapObject) Empty() bool {
	return ego.ptr.Count() == 0
}

/*
Checks if the content of the object is equal to the content of another object.
Nested objects and lists are compared recursively (by value).

Parameters:
  - another - an object to compare with.

Returns:
  - true if the objects are equal, false otherwise.
*/
func (ego *MapObject) Equals(another Object) bool {
	ego.assert()
	return ego.ptr.isEqual(another)
}

/*
Creates a new object containing all elements of the old object and another object.
The old object remains unchanged.
If both objects contain a key, the value from another object is used.

Parameters:
  - another - an object to merge.

Returns:
  - new object.
*/
func (ego *MapObject) Merge(another Object) Object {
	ego.assert()
	result := ego.Clone()
	another.ForEach(func(key string, val any) {
		result.Set(key, val)
	})
	return result
}

/*
Creates a new object containing the given fields of the existing object.

Parameters:
  - keys... - any amount of keys to be in the new object.

Returns:
  - created plucked object.
*/
func (ego *MapObject) Pluck(keys ...string) Object {
	ego.assert()
	result := NewObject()
	for _, key := range keys {
		result.Set(key, ego.Get(key))
	}
	return result
}

/*
Checks if the object contains a field with a given value.
Objects and lists are compared by reference.

Parameters:
  - value - the value to check.

Returns:
  - true if the object contains the value, false otherwise.
*/
func (ego *MapObject) Contains(value any) bool {
	ego.assert()
	for _, item := range ego.val {
		switch item.(type) {
		case Object, List:
			if item == value {
				return true
			}
		default:
			if item.getVal() == value {
				return true
			}
		}
	}
	return false
}

/*
Gives a key containing a given value.
If multiple keys contain the value, any of them is returned.

Parameters:
  - value - the value to check.

Returns:
  - key for the value (empty string if the object does not contain the value).
*/
func (ego *MapObject) KeyOf(value any) string {
	ego.assert()
	for key, item := range ego.val {
		switch item.(type) {
		case Object, List:
			if item == value {
				return key
			}
		default:
			if item.getVal() == value {
				return key
			}
		}
	}
	return ""
}

/*
Checks if a given key exists within the object.

Parameters:
  - key - the key to check.

Returns:
  - true if the key exists, false otherwise.
*/
func (ego *MapObject) KeyExists(key string) bool {
	ego.assert()
	_, ok := ego.val[key]
	return ok
}

/*
Executes a given function over an every field of the object.
The function has two parameters: key of the current field and its value.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEach(function func(string, any)) Object {
	ego.assert()
	for key, item := range ego.val {
		function(key, item.getVal())
	}
	return ego.ptr
}

/*
Executes a given function over an every field of the object.
The function has one parameter, value of the current field.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachValue(function func(any)) Object {
	ego.assert()
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.ptr
}

/*
Executes a given function over all objects nested in the object.
Fields with other types are ignored.
The function has one parameter, the current object.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachObject(function func(Object)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(Object)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all lists in the object.
Fields with other types are ignored.
The function has one parameter, the current list.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachList(function func(List)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(List)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all strings in the object.
Fields with other types are ignored.
The function has one parameter, the current string.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachString(function func(string)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all bools in the object.
Fields with other types are ignored.
The function has one parameter, the current bool.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachBool(function func(bool)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all ints in the object.
Fields with other types are ignored.
The function has one parameter, the current int.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachInt(function func(int)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all floats in the object.
Fields with other types are ignored.
The function has one parameter, the current float.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachFloat(function func(float64)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Copies the object and modifies each field by a given mapping function.
The resulting field can have a different type than the original one.
The function has two parameters: current key and value of the current element. Returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *MapObject) Map(function func(string, any) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(key, item.getVal()))
	}
	return result
}

/*
Copies the object and modifies each field by a given mapping function.
The resulting field can have a different type than the original one.
The function has one parameter, value of the current field, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapValues(function func(any) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(item.getVal()))
	}
	return result
}

/*
Selects all nested objects of the object and modifies each of them by a given mapping function.
Fields with other types are ignored.
The resulting field can have a different type than the original one.
The function has one parameter, the current object, and returns empty interface.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapObjects(function func(Object) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		val, ok := item.(Object)
		if ok {
			result.Set(key, function(val))
		}
	}
	return result
}

/*
Selects all lists of the object and modifies each of them by a given mapping function.
Fields with other types are ignored.
The resulting field can have a different type than the original one.
The function has one parameter, the current list, and returns empty interface.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapLists(function func(List) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		val, ok := item.(List)
		if ok {
			result.Set(key, function(val))
		}
	}
	return result
}

/*
Selects all nested strings of the object and modifies each of them by a given mapping function.
Fields with other types are ignored.
The resulting field can have a different type than the original one.
The function has one parameter, the current string, and returns empty interface.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapStrings(function func(string) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			result.Set(key, function(val))
		}
	}
	return result
}

/*
Selects all nested ints of the object and modifies each of them by a given mapping function.
Fields with other types are ignored.
The resulting field can have a different type than the original one.
The function has one parameter, the current int, and returns empty interface.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapInts(function func(int) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result.Set(key, function(val))
		}
	}
	return result
}

/*
Selects all nested floats of the object and modifies each of them by a given mapping function.
Fields with other types are ignored.
The resulting field can have a different type than the original one.
The function has one parameter, the current float, and returns empty interface.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapFloats(function func(float64) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			result.Set(key, function(val))
		}
	}
	return result
}

/*
Parallelly executes a given function over an every field of the object.
The function has two parameters: key of the current field and its value.
The order of the iterations is random.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged object.
*/
func (ego *MapObject) ForEachAsync(function func(string, any)) Object {
	ego.assert()
	var wg sync.WaitGroup
	step := func(group *sync.WaitGroup, k string, x any) {
		function(k, x)
		group.Done()
	}
	wg.Add(ego.Count())
	for key, item := range ego.val {
		go step(&wg, key, item.getVal())
	}
	wg.Wait()
	return ego.ptr
}

/*
Copies the object and paralelly modifies each field by a given mapping function.
The resulting field can have a different type than the original one.
The function has two parameters: key of the current field and its value.
The old object remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new object.
*/
func (ego *MapObject) MapAsync(function func(string, any) any) Object {
	ego.assert()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(ego.Count())
	result := NewObject()
	step := func(group *sync.WaitGroup, k string, x any) {
		mutex.Lock()
		result.Set(k, function(k, x))
		mutex.Unlock()
		group.Done()
	}
	for key, item := range ego.val {
		go step(&wg, key, item.getVal())
	}
	wg.Wait()
	return result
}

/*
Acquires the element specified by the given tree form.

Parameters:
  - tf - tree form string.

Returns:
  - corresponding value (any type, has to be asserted).
*/
func (ego *MapObject) GetTF(tf string) any {
	ego.assert()
	if tf[0] != '.' || len(tf) < 2 {
		panic("'" + tf + "' is not a valid tree form.")
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		return ego.ptr.GetObject(tf[:dot]).GetTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		return ego.ptr.GetList(tf[:hash]).GetTF(tf[hash:])
	}
	return ego.ptr.Get(tf)
}

/*
Sets the element specified by the given tree form.

Parameters:
  - tf - tree form string,
  - value - value to set.

Returns:
  - updated object.
*/
func (ego *MapObject) SetTF(tf string, value any) Object {
	ego.assert()
	if tf[0] != '.' || len(tf) < 2 {
		panic("'" + tf + "' is not a valid tree form.")
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		key := tf[:dot]
		var object Object
		if ego.KeyExists(key) {
			object = ego.GetObject(key)
		} else {
			object = NewObject()
			ego.Set(key, object)
		}
		object.SetTF(tf[dot:], value)
		return ego.ptr
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		key := tf[:hash]
		var list List
		if ego.KeyExists(key) {
			list = ego.GetList(key)
		} else {
			list = NewList()
			ego.Set(key, list)
		}
		list.SetTF(tf[hash:], value)
		return ego.ptr
	}
	return ego.ptr.Set(tf, value)
}
