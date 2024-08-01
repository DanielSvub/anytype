/*
AnyType Library for Go
Object (dictionary) implementation
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
Map object, a reference type. Contains a map.

Implements:
  - field,
  - Object.
*/
type object struct {
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
func NewObject(values ...any) Object {
	ego := &object{val: make(map[string]field)}
	ego.ptr = ego
	ego.Set(values...)
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
		panic("unsupported map type")
	}
	return object
}

/*
Asserts that the object is initialized.
*/
func (ego *object) assert() {
	if ego == nil || ego.val == nil {
		panic("object is not initialized")
	}
}

/*
Defined in the field interface.
Acquires the value of the field, in this case a reference to the whole struct (Object is a reference type).

Returns:
  - value of the field.
*/
func (ego *object) getVal() any {
	return ego.ptr
}

/*
Defined in the field interface.
Creates a deep copy of the field, in this case a new object with identical fields.
Can be called recursively.

Returns:
  - deep copy of the field.
*/
func (ego *object) copy() any {
	obj := NewObject()
	for key, value := range ego.val {
		obj.Set(key, value.copy())
	}
	return obj
}

/*
Defined in the field interface.
Serializes the field into the JSON format, in this case prints all keys and their values.
Can be called recursively.

Returns:
  - string representing serialized field.
*/
func (ego *object) serialize() string {
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
Defined in the field interface.
Checks if the content of the field is equal to the given field.
Can be called recursively.

Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *object) isEqual(another any) bool {
	obj, ok := another.(*object)
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

func (ego *object) Init(ptr Object) {
	ego.ptr = ptr
}

func (ego *object) Ego() Object {
	return ego.ptr
}

func (ego *object) Set(values ...any) Object {
	ego.assert()
	length := len(values)
	if length%2 != 0 {
		panic("object fields have to be set as key-value pairs")
	}
	for i := 0; i < length; i += 2 {
		name, ok := values[i].(string)
		if !ok {
			panic("object key has to be string")
		}
		ego.val[name] = parseVal(values[i+1])
	}
	return ego.ptr
}

func (ego *object) Unset(keys ...string) Object {
	ego.assert()
	for _, key := range keys {
		if ego.val[key] == nil {
			panic(fmt.Sprintf("object does not have a field '%s'", key))
		}
		delete(ego.val, key)
	}
	return ego.ptr
}

func (ego *object) Clear() Object {
	ego.assert()
	ego.val = make(map[string]field, 0)
	return ego.ptr
}

func (ego *object) Get(key string) any {
	ego.assert()
	if ego.val[key] == nil {
		panic(fmt.Sprintf("object does not have a field '%s'", key))
	}
	obj := ego.val[key]
	switch obj.(type) {
	case Object, List:
		return obj
	default:
		return obj.getVal()
	}
}

func (ego *object) GetObject(key string) Object {
	o, ok := ego.Get(key).(Object)
	if !ok {
		panic("field is not an object")
	}
	return o
}

func (ego *object) GetList(key string) List {
	o, ok := ego.Get(key).(List)
	if !ok {
		panic("field is not a list")
	}
	return o
}

func (ego *object) GetString(key string) string {
	o, ok := ego.Get(key).(string)
	if !ok {
		panic("field is not a string")
	}
	return o
}

func (ego *object) GetBool(key string) bool {
	o, ok := ego.Get(key).(bool)
	if !ok {
		panic("field is not a bool")
	}
	return o
}

func (ego *object) GetInt(key string) int {
	o, ok := ego.Get(key).(int)
	if !ok {
		panic("field is not an int")
	}
	return o
}

func (ego *object) GetFloat(key string) float64 {
	o, ok := ego.Get(key).(float64)
	if !ok {
		panic("field is not a float")
	}
	return o
}

func (ego *object) TypeOf(key string) Type {
	ego.assert()
	if !ego.KeyExists(key) {
		return TypeUndefined
	}
	switch ego.val[key].(type) {
	case *atNil:
		return TypeNil
	case *atString:
		return TypeString
	case *atInt:
		return TypeInt
	case *atBool:
		return TypeBool
	case *atFloat:
		return TypeFloat
	case *object:
		return TypeObject
	case *list:
		return TypeList
	default:
		panic("unknown field type")
	}
}

func (ego *object) String() string {
	ego.assert()
	return ego.ptr.serialize()
}

func (ego *object) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic("invalid indentation")
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

func (ego *object) Dict() map[string]any {
	ego.assert()
	dict := make(map[string]any, 0)
	for key, value := range ego.val {
		dict[key] = value.getVal()
	}
	return dict
}

func (ego *object) Keys() List {
	keys := NewList()
	for key := range ego.val {
		keys.Add(key)
	}
	return keys
}

func (ego *object) Values() List {
	values := NewList()
	for _, value := range ego.val {
		values.Add(value.getVal())
	}
	return values
}

func (ego *object) Clone() Object {
	ego.assert()
	return ego.ptr.copy().(*object)
}

func (ego *object) Count() int {
	ego.assert()
	return len(ego.val)
}

func (ego *object) Empty() bool {
	return ego.ptr.Count() == 0
}

func (ego *object) Equals(another Object) bool {
	ego.assert()
	return ego.ptr.isEqual(another)
}

func (ego *object) Merge(another Object) Object {
	ego.assert()
	result := ego.Clone()
	another.ForEach(func(key string, val any) {
		result.Set(key, val)
	})
	return result
}

func (ego *object) Pluck(keys ...string) Object {
	ego.assert()
	result := NewObject()
	for _, key := range keys {
		result.Set(key, ego.Get(key))
	}
	return result
}

func (ego *object) Contains(value any) bool {
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

func (ego *object) KeyOf(value any) string {
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
	panic(fmt.Sprintf("object does not contain value %v", value))
}

func (ego *object) KeyExists(key string) bool {
	ego.assert()
	_, ok := ego.val[key]
	return ok
}

func (ego *object) ForEach(function func(string, any)) Object {
	ego.assert()
	for key, item := range ego.val {
		function(key, item.getVal())
	}
	return ego.ptr
}

func (ego *object) ForEachValue(function func(any)) Object {
	ego.assert()
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.ptr
}

func (ego *object) ForEachObject(function func(Object)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(Object)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) ForEachList(function func(List)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(List)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) ForEachString(function func(string)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) ForEachBool(function func(bool)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) ForEachInt(function func(int)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) ForEachFloat(function func(float64)) Object {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *object) Map(function func(string, any) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(key, item.getVal()))
	}
	return result
}

func (ego *object) MapValues(function func(any) any) Object {
	ego.assert()
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(item.getVal()))
	}
	return result
}

func (ego *object) MapObjects(function func(Object) any) Object {
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

func (ego *object) MapLists(function func(List) any) Object {
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

func (ego *object) MapStrings(function func(string) any) Object {
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

func (ego *object) MapInts(function func(int) any) Object {
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

func (ego *object) MapFloats(function func(float64) any) Object {
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

func (ego *object) ForEachAsync(function func(string, any)) Object {
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

func (ego *object) MapAsync(function func(string, any) any) Object {
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

func (ego *object) GetTF(tf string) any {
	ego.assert()
	if tf[0] != '.' || len(tf) < 2 {
		panic(fmt.Sprintf("'%s' is not a valid tree form", tf))
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		key := tf[:dot]
		if !ego.ptr.KeyExists(key) || ego.ptr.TypeOf(key) != TypeObject {
			return nil
		}
		return ego.ptr.GetObject(key).GetTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		key := tf[:hash]
		if !ego.ptr.KeyExists(key) || ego.ptr.TypeOf(key) != TypeList {
			return nil
		}
		return ego.ptr.GetList(key).GetTF(tf[hash:])
	}
	if !ego.ptr.KeyExists(tf) {
		return nil
	}
	return ego.ptr.Get(tf)
}

func (ego *object) SetTF(tf string, value any) Object {
	ego.assert()
	if tf[0] != '.' || len(tf) < 2 {
		panic(fmt.Sprintf("'%s' is not a valid tree form", tf))
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
