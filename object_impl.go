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
	ego := &object{val: map[string]field{}}
	ego.Init(ego)
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
Defined in the field interface.
Acquires the value of the field, in this case a reference to the whole struct (Object is a reference type).

Returns:
  - value of the field.
*/
func (ego *object) getVal() any {
	return ego.Ego()
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
	var result strings.Builder
	result.WriteRune('{')
	i := 0
	for field, value := range ego.val {
		result.WriteString(fmt.Sprintf("%s:%s", strconv.Quote(field), value.serialize()))
		if i++; i < len(ego.val) {
			result.WriteRune(',')
		}
	}
	result.WriteRune('}')
	return result.String()
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
	if !ok || ego.Ego().Count() != obj.Count() {
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
	length := len(values)
	if length&1 == 1 {
		panic("object fields have to be set as key-value pairs")
	}
	for i := 0; i < length; i += 2 {
		name, ok := values[i].(string)
		if !ok {
			panic("object key has to be string")
		}
		ego.val[name] = parseVal(values[i+1])
	}
	return ego.Ego()
}

func (ego *object) Unset(keys ...string) Object {
	for _, key := range keys {
		delete(ego.val, key)
	}
	return ego.Ego()
}

func (ego *object) Clear() Object {
	ego.val = map[string]field{}
	return ego.Ego()
}

func (ego *object) Get(key string) any {
	field, exists := ego.val[key]
	if !exists {
		panic(fmt.Sprintf("object does not have a field '%s'", key))
	}
	return field.getVal()

}

func (ego *object) GetObject(key string) Object {
	o, ok := ego.Get(key).(Object)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not an object", key))
	}
	return o
}

func (ego *object) GetList(key string) List {
	o, ok := ego.Get(key).(List)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not a list", key))
	}
	return o
}

func (ego *object) GetString(key string) string {
	o, ok := ego.Get(key).(string)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not a string", key))
	}
	return o
}

func (ego *object) GetBool(key string) bool {
	o, ok := ego.Get(key).(bool)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not a bool", key))
	}
	return o
}

func (ego *object) GetInt(key string) int {
	o, ok := ego.Get(key).(int)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not an int", key))
	}
	return o
}

func (ego *object) GetFloat(key string) float64 {
	o, ok := ego.Get(key).(float64)
	if !ok {
		panic(fmt.Sprintf("field '%s' is not a float", key))
	}
	return o
}

func (ego *object) TypeOf(key string) Type {
	switch ego.val[key].(type) {
	case Object:
		return TypeObject
	case List:
		return TypeList
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
	default:
		return TypeUndefined
	}
}

func (ego *object) String() string {
	return ego.Ego().serialize()
}

func (ego *object) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic(fmt.Sprintf("indentation %d is not between 1 and 10", indent))
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

func (ego *object) Dict() map[string]any {
	dict := make(map[string]any, ego.Ego().Count())
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
	return ego.Ego().copy().(*object)
}

func (ego *object) Count() int {
	return len(ego.val)
}

func (ego *object) Empty() bool {
	return ego.Ego().Count() == 0
}

func (ego *object) Equals(another Object) bool {
	return ego.Ego().isEqual(another)
}

func (ego *object) Merge(another Object) Object {
	result := ego.Clone()
	another.ForEach(func(key string, val any) {
		result.Set(key, val)
	})
	return result
}

func (ego *object) Pluck(keys ...string) Object {
	result := NewObject()
	for _, key := range keys {
		result.Set(key, ego.Get(key))
	}
	return result
}

func (ego *object) Contains(value any) bool {
	for _, item := range ego.val {
		if item.getVal() == value {
			return true
		}
	}
	return false
}

func (ego *object) KeyOf(value any) string {
	for key, item := range ego.val {
		if item.getVal() == value {
			return key
		}
	}
	panic(fmt.Sprintf("object does not contain value %v", value))
}

func (ego *object) KeyExists(key string) bool {
	_, ok := ego.val[key]
	return ok
}

func (ego *object) ForEach(function func(string, any)) Object {
	for key, item := range ego.val {
		function(key, item.getVal())
	}
	return ego.Ego()
}

func (ego *object) ForEachValue(function func(any)) Object {
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.Ego()
}

func (ego *object) ForEachObject(function func(Object)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(Object)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) ForEachList(function func(List)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(List)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) ForEachString(function func(string)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) ForEachBool(function func(bool)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) ForEachInt(function func(int)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) ForEachFloat(function func(float64)) Object {
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *object) Map(function func(string, any) any) Object {
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(key, item.getVal()))
	}
	return result
}

func (ego *object) MapValues(function func(any) any) Object {
	result := NewObject()
	for key, item := range ego.val {
		result.Set(key, function(item.getVal()))
	}
	return result
}

func (ego *object) MapObjects(function func(Object) any) Object {
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
	return ego.Ego()
}

func (ego *object) MapAsync(function func(string, any) any) Object {
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
	if len(tf) < 2 || tf[0] != '.' {
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
		return ego.Ego().GetObject(key).GetTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		key := tf[:hash]
		if !ego.ptr.KeyExists(key) || ego.ptr.TypeOf(key) != TypeList {
			return nil
		}
		return ego.Ego().GetList(key).GetTF(tf[hash:])
	}
	if !ego.ptr.KeyExists(tf) {
		return nil
	}
	return ego.Ego().Get(tf)
}

func (ego *object) SetTF(tf string, value any) Object {
	if len(tf) < 2 || tf[0] != '.' {
		panic(fmt.Sprintf("'%s' is not a valid tree form", tf))
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		key := tf[:dot]
		var object Object
		if ego.TypeOf(key) == TypeObject {
			object = ego.GetObject(key)
		} else {
			object = NewObject()
			ego.Set(key, object)
		}
		object.SetTF(tf[dot:], value)
		return ego.Ego()
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
		return ego.Ego()
	}
	return ego.Ego().Set(tf, value)
}
